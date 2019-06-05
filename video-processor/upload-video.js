const exec = require('util').promisify(require('child_process').exec);
const fs   = require('fs-extra');
const path = require('path');
const tmp  = require('tmp-promise');
const B2   = require('./b2');

function loadEnv () {
	const res = require('dotenv').config();
	if (res.error) {
		throw res.error;
	}
}

async function initB2 () {
	await B2.init({
		accountId:      process.env.B2_ACCOUNT_ID,
		applicationKey: process.env.B2_APPLICATION_KEY,
		bucket:         'videoso',
	});
}

async function processVideo (videoPath, dashDir) {

	// Transcode
	const localVideoPath = `${dashDir}/input.mp4`;
	await fs.symlink(videoPath, localVideoPath);
	await exec(`./transcode-video.sh "${localVideoPath}"`);

	// Rename transcoded videos
	const rawTranscodedVideoPaths = [
		`${dashDir}/input_video_360.mp4`,
		`${dashDir}/input_video_480.mp4`,
		`${dashDir}/input_video_720.mp4`,
		`${dashDir}/input_video_1080.mp4`,
		`${dashDir}/input_audio_192.m4a`,
	];
	const transcodedVideoPaths = [
		`${dashDir}/video_360.mp4`,
		`${dashDir}/video_480.mp4`,
		`${dashDir}/video_720.mp4`,
		`${dashDir}/video_1080.mp4`,
		`${dashDir}/audio_192.m4a`,
	];
	for (let i = 0; i < transcodedVideoPaths.length; ++i) {
		const p1 = rawTranscodedVideoPaths[i];
		const p2 = transcodedVideoPaths[i];
		await fs.move(p1, p2);
	}

	// Generate MPD
	const mpdPath = `${dashDir}/playlist.mpd`;
	const mpdCmd = [
		'./generate-mpd.sh',
		`"${mpdPath}"`,
		`"${transcodedVideoPaths.join('" "')}"`,
	].join(' ');
	await exec(mpdCmd);

	// Return the important file paths
	const dashVideoPaths = transcodedVideoPaths.map(path => {
		if (path.includes('.mp4')) {
			return path.replace('.mp4', '_dashinit.mp4');
		} else {
			return path.replace('.m4a', '_dashinit.mp4');
		}
	});
	return dashVideoPaths.concat(mpdPath);
}

async function main () {
	console.log('Initializing...');
	loadEnv();
	await initB2();
	console.log('Done.');

	const args = process.argv.slice(2);
	const videoPath = path.resolve(args[0]);
	const videoBase = args[1] || path.parse(videoPath).base;
	const videoName = path.parse(videoBase).name;

	console.log('Uploading raw video...');
	const rawVideoUrl = await B2.uploadFile(videoPath, `videos-raw/${videoBase}`);
	console.log('Raw video available at:', rawVideoUrl);

	console.log('Processing video for MPEG-DASH...');
	const dir = await tmp.dir();
	const files = await processVideo(videoPath, dir.path);
	console.log('Done.');

	console.log('Uploading dash files...');
	let dashVideoUrl;
	await Promise.all(
		files.map(async fsPath => {
			const dashBase = fsPath.substring(dir.path.length + 1);
			const b2Path = `videos-dash/${videoName}/${dashBase}`;
			const b2Url = await B2.uploadFile(fsPath, b2Path);
			if (dashBase.includes('.mpd')) {
				dashVideoUrl = b2Url;
			}
		})
	);
	console.log('Dash video available at:', dashVideoUrl);
}

main().catch(err => {
	console.error('error in main:', err);
});
