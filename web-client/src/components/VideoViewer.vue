<template>
	<div class="video-viewer">
		<div class="video-wrapper">
			<video
				playsinline
				controls
				ref="video">
				<source :src="video.videoUrl">
			</video>
		</div>
		<div class="primary-info">
			<h1>{{ video.title }}</h1>
			<div class="view-count">
				{{ formatViewCount(video.viewCount) }}
			</div>
		</div>
		<div class="secondary-info">
			<div class="publish-info">
				Published by
				<router-link :to="`/channel/${video.author.name}`">
					{{ video.author.name }}
				</router-link>
				on
				{{ formatDate(video.publishedAt) }}
			</div>
			<div class="description" v-if="video.description">
				{{ video.description }}
			</div>
			<div class="description empty-description" v-else>
				No description.
			</div>
		</div>
	</div>
</template>

<script>
	import Plyr from 'plyr';

	export default {
		name: 'VideoViewer',

		props: {
			/**
			 * @typedef {object} Video
			 * @property {string} title
			 * @property {string} description
			 * @property {number} viewCount
			 * @property {string} videoUrl
			 * @property {object} author
			 * @property {number} author.id
			 * @property {string} author.name
			 * @property {Date}   publishedAt
			 */
			video: Object,
		},

		mounted () {
			new Plyr(this.$refs.video, {
				controls: [
					'play-large',
					'play',
					'progress',
					'current-time',
					'mute',
					'volume',
					'captions',
					'settings',
					'airplay',
					'fullscreen',
				]
			});
		},

		methods: {
			formatDate (date) {
				return date.toLocaleString('en-US', {
					month: 'short',
					day: 'numeric',
					year: 'numeric',
				});
			},

			formatViewCount (viewCount) {
				if (viewCount === 0) {
					return 'No views';
				}
				if (viewCount === 1) {
					return '1 view';
				}
				return `${viewCount} views`;
			},
		},
	}
</script>

<style lang="scss" scoped>
	.video-wrapper {
		background-color: black;
		position: relative;

		video {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			object-fit: contain;
		}
	}

	.primary-info {
		display: flex;
		margin: 2rem 0;
		//padding: 0 1rem;

		h1 {
			flex: 1;
			font-size: 3.2rem;
			margin: 0;
		}

		.view-count {
			font-size: 2.0rem;
		}
	}

	.secondary-info {
		font-size: 1.4rem;
		//padding: 0 1rem;

		.publish-info {
			a {
				font-weight: bold;
			}
		}

		.description {
			margin: 2rem 0 0;
		}

		.empty-description {
			font-style: italic;
		}
	}
</style>
