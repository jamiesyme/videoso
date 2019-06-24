const fs = require('fs-extra');
const B2 = require('backblaze-b2');

let client;
let bucketName;
let bucketId;

/**
 * @param {object} params
 * @param {string} params.accountId
 * @param {string} params.applicationKey
 * @param {string} params.bucket
 */
async function init (params) {
	client = new B2({
		accountId: params.accountId,
		applicationKey: params.applicationKey,
	});
	await client.authorize();

	bucketName = params.bucket;
	const res = await client.getBucket({ bucketName });
	if (res.data &&
	    res.data.buckets &&
	    res.data.buckets.length > 0 &&
	    res.data.buckets[0].bucketId) {
		bucketId = res.data.buckets[0].bucketId;
	} else {
		throw new Error('failed to get bucket id');
	}
}

/**
 * @param {string} fsPath
 * @param {string} b2Path
 */
async function uploadFile (fsPath, b2Path) {
	const fileBuffer = await fs.readFile(fsPath);
	const uploadUrlRes = await client.getUploadUrl({ bucketId });
	try {
		await client.uploadFile({
			uploadUrl:       uploadUrlRes.data.uploadUrl,
			uploadAuthToken: uploadUrlRes.data.authorizationToken,
			fileName:        b2Path,
			data:            fileBuffer,
		});
		return `${client.downloadUrl}/file/${bucketName}/${b2Path}`;
	} catch (err) {
		if (err.response.status === 503) {
			return await uploadFile(fsPath, b2Path);
		} else {
			throw err;
		}
	}
}

module.exports = {
	init,
	uploadFile,
};
