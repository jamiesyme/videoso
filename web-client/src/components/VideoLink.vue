<template>
	<router-link
		class="video-link"
		:to="`/video/${video.id}`">
		<div
			class="thumbnail-wrapper aspect-16-9"
			:style="extraStyle">
			<img :src="video.thumbnailUrl">
			<div class="duration">
				{{ formatDuration(video.duration) }}
			</div>
		</div>
		<div class="video-title">{{ video.title }}</div>
		<div class="video-publisher">
			{{ video.author.name }}
		</div>
	</router-link>
</template>

<script>
	export default {
		name: 'VideoLink',

		props: {
			/**
			 * @typedef {object} Video
			 * @property {number} id
			 * @property {string} title
			 * @property {string} thumbnailUrl
			 * @property {object} author
			 * @property {number} author.id
			 * @property {string} author.name
			 */
			video: Object,

			/**
			 * Use on mobile to stretch thumbnail to edges of screen.
			 * @param {string} extraWidth - '10px', '2rem', etc.
			 */
			extraWidth: String,
		},

		methods: {
			formatDuration (seconds) {
				function pad2 (num) {
					return num.toString().padStart(2, '0');
				}
				const secs = Math.floor(seconds % 60);
				const mins = Math.floor(seconds / 60);
				const hrs = Math.floor(seconds / 60 / 60);
				if (hrs > 0) {
					return `${pad2(hrs)}:${pad2(mins)}:${pad2(secs)}`;
				} else {
					return `${pad2(mins)}:${pad2(secs)}`;
				}
			},
		},

		computed: {
			extraStyle () {
				if (this.extraWidth) {
					return {
						marginLeft: `calc(-0.5 * ${this.extraWidth})`,
						width: `calc(100% + ${this.extraWidth})`,
					};
				}
				return null;
			},
		},
	}
</script>

<style lang="scss" scoped>
	.video-link {
		display: block;
	}

	.thumbnail-wrapper {
		background-color: black;
		position: relative;

		img {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			object-fit: contain;
		}

		.duration {
			background-color: rgba(0, 0, 0, 0.8);
			color: rgba(255, 255, 255, 0.9);
			position: absolute;
			bottom: 0.4rem;
			right: 0.4rem;
			padding: 0 0.4rem;
			font-size: 1.3rem;
			font-weight: bold;
		}
	}

	.video-title {
		font-size: 1.4rem;
		font-weight: bold;
		margin: 0.8rem 0 0;
	}

	.video-publisher {
		color: #606c76;
		font-size: 1.3rem;
	}
</style>
