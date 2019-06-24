<template>
	<div class="video-viewer">
		<div
			class="video-wrapper"
			:style="videoStyle">
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
			<div
				class="description"
				v-html="videoDescriptionHtml"
				v-if="video.description">
			</div>
			<div
				class="description empty-description"
				v-else>
				No description.
			</div>
		</div>
	</div>
</template>

<script>
	import DashJs from 'dashjs';
	import Plyr from 'plyr';

	function sanitizeHtml (text) {
		const temp = document.createElement('div');
		temp.textContent = text;
		return temp.innerHTML;
	}

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

			/**
			 * Use on mobile to stretch video to edges of screen.
			 * @param {string} extraWidth - '10px', '2rem', etc.
			 */
			extraWidth: String,
		},

		data () {
			return {
				dash: null,
				player: null,
			};
		},

		mounted () {
			this.dash = DashJs.MediaPlayer().create();
			this.dash.initialize(null, null, false);
			this.player = new Plyr(this.$refs.video, {
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
			this.updateVideoElement();
		},

		methods: {
			updateVideoElement () {
				this.dash.reset();
				const isDashUrl = /\.mpd$/.test(this.video.videoUrl);
				if (isDashUrl) {
					this.dash.attachView(this.$refs.video);
					this.dash.attachSource(this.video.videoUrl);
				}
			},

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

		computed: {
			videoStyle () {
				if (this.extraWidth) {
					return {
						marginLeft: `calc(-0.5 * ${this.extraWidth})`,
						width: `calc(100% + ${this.extraWidth})`,
						borderTop: '0.1rem solid #eee',
						borderBottom: '0.1rem solid #eee',
					};
				}
				return null;
			},

			videoDescriptionHtml () {
				return this.video.description.split('\n\n').map(p => {
					const safeP = sanitizeHtml(p);
					const fullP = safeP.trim().replace(/\n/g, '<br>');
					return `<p>${fullP}</p>`;
				}).join('');
			},
		},

		watch: {
			'video.videoUrl': function () {
				this.updateVideoElement();
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
			object-fit: cover;
		}
	}

	.primary-info {
		display: flex;
		margin: 2rem 0;

		h1 {
			flex: 1;
			font-size: 3.2rem;
			margin: 0;
		}

		.view-count {
			font-size: 2.0rem;
			margin-left: 2rem;
		}
	}

	.secondary-info {
		font-size: 1.4rem;

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

<style lang="scss">
	// Hack to style the paragraphs that are dynamically inserted
	.video-viewer .description p {
		margin-bottom: 1.5rem;
	}
</style>
