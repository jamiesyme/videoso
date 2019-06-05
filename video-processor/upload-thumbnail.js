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

async function processImage (imagePath, outputDir) {

	// Transcode
	const localImagePath = `${outputDir}/image.jpg`;
	await fs.symlink(imagePath, localImagePath);
	await exec(`./transcode-image.sh "${localImagePath}"`);

	// Return the new file paths
	return [
		`${outputDir}/image_240.jpg`,
		`${outputDir}/image_360.jpg`,
		`${outputDir}/image_480.jpg`,
		`${outputDir}/image_720.jpg`,
		`${outputDir}/image_1080.jpg`,
	];
}

async function main () {
	console.log('Initializing...');
	loadEnv();
	await initB2();
	console.log('Done.');

	const args = process.argv.slice(2);
	const imagePath = path.resolve(args[0]);
	const imageBase = args[1] || path.parse(imagePath).base;
	const imageName = path.parse(imageBase).name;

	console.log('Uploading raw image...');
	const rawImageUrl = await B2.uploadFile(imagePath, `thumbnails-raw/${imageBase}`);
	console.log('Raw image available at:', rawImageUrl);

	console.log('Transcoding image...');
	const dir = await tmp.dir();
	const files = await processImage(imagePath, dir.path);
	console.log('Done.');

	console.log('Uploading transcoded images...');
	let bigImageUrl;
	await Promise.all(
		files.map(async fsPath => {
			const thumbBase = fsPath.substring(dir.path.length + 1);
			const b2Path = `thumbnails-converted/${imageName}/${thumbBase}`;
			const b2Url = await B2.uploadFile(fsPath, b2Path);
			if (thumbBase.includes('_1080.jpg')) {
				bigImageUrl = b2Url;
			}
		})
	);
	console.log('Transcoded image(s) available at:', bigImageUrl);
}

main().catch(err => {
	console.error('error in main:', err);
});
