const exec = require('util').promisify(require('child_process').exec);
const fs   = require('fs-extra');
const path = require('path');
const tmp  = require('tmp-promise');
const B2   = require('backblaze-b2');

let b2Client;
let b2BucketName;
let b2BucketId;
const tmpDir = path.resolve('./tmp');

function loadEnv () {
	const res = require('dotenv').config();
	if (res.error) {
		throw res.error;
	}
}

async function initB2 () {
	b2Client = new B2({
		accountId: process.env.B2_ACCOUNT_ID,
		applicationKey: process.env.B2_APPLICATION_KEY,
	});
	await b2Client.authorize();

	b2BucketName = 'videoso';
	const res = await b2Client.getBucket({ bucketName: b2BucketName });
	if (res.data &&
	    res.data.buckets &&
	    res.data.buckets.length > 0 &&
	    res.data.buckets[0].bucketId) {
		b2BucketId = res.data.buckets[0].bucketId;
	} else {
		throw new Error('failed to get bucket id');
	}
}

async function uploadFile (fsPath, b2Path) {
	const fileBuffer = await fs.readFile(fsPath);
	const uploadUrlRes = await b2Client.getUploadUrl({
		bucketId: b2BucketId,
	});
	try {
		await b2Client.uploadFile({
			uploadUrl:       uploadUrlRes.data.uploadUrl,
			uploadAuthToken: uploadUrlRes.data.authorizationToken,
			fileName:        b2Path,
			data:            fileBuffer,
		});
		return `${b2Client.downloadUrl}/file/${b2BucketName}/${b2Path}`;
	} catch (err) {
		if (err.response.status === 503) {
			return await uploadFile(fsPath, b2Path);
		} else {
			throw err;
		}
	}
}

async function processVideo (videoPath, dashDir) {

	// Transcode
	const localVideoPath = `${dashDir}/input.mp4`;
	await fs.symlink(videoPath, localVideoPath);
	await exec(`./transcode-video.sh "${localVideoPath}"`);

	// Generate MPD
	const mpdPath = `${dashDir}/playlist.mpd`;
	const transcodedVideoPaths = [
		`${dashDir}/input_video_360_28.mp4`,
		`${dashDir}/input_video_480_23.mp4`,
		`${dashDir}/input_video_720_20.mp4`,
		`${dashDir}/input_video_1080_18.mp4`,
		`${dashDir}/input_audio_192.m4a`,
	];
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
	const rawVideoUrl = await uploadFile(videoPath, `videos-raw/${videoBase}`);
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
			const b2Url = await uploadFile(fsPath, b2Path);
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
