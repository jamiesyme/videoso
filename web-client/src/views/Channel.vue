<template>
	<div class="channel-view">
		<section
			class="banner-section"
			v-if="user">
			<img :src="user.bannerUrl">
			<div class="container">
				<h1>{{ user.name }}'s Channel</h1>
			</div>
		</section>
		<div class="container" v-if="video">
			<section class="video-section">
				<VideoViewer
					:video="video"
					v-if="video">
				</VideoViewer>
				<hr>
			</section>
			<section class="related-section">
				<h2>Other Videos</h2>
				<VideoLink
					:video="video"
					:key="video.id"
					v-for="video in otherVideos">
				</VideoLink>
				<hr>
			</section>
		</div>

		<div class="container" v-else>
			No videos.
		</div>
	</div>
</template>

<script>
	import Content from '@/content';
	import ContentUtils from '@/utils/content';
	import VideoLink from '@/components/VideoLink';
	import VideoViewer from '@/components/VideoViewer';

	export default {
		components: {
			VideoLink,
			VideoViewer,
		},

		data () {
			return {
				user: null,
				video: null,
				otherVideos: [],
			};
		},

		mounted () {
			this.refreshContent();
		},

		methods: {
			async refreshContent () {
				const dummyBannerUrl = 'https://dummyimage.com/1920x200/eee/e6e6e6';
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				// Load content (if hasn't already been loaded)
				if (Content.users.length === 0) {
					await Content.load();
				}

				// Find the user
				this.user = Content.users.find(user => {
					function tr (str) {
						return str.trim().toLowerCase();
					}
					const routeName = tr(this.$route.params.userName);
					const userName = tr(user.name);
					return routeName === userName;
				});
				if (!this.user) {
					this.videos = null;
					this.otherVideos = [];
					return;
				} else {
					if (!this.user.bannerUrl) {
						this.user.bannerUrl = dummyBannerUrl;
					}
				}

				// Find the user's videos
				const userVideos = Content.videos.filter(vid => {
					return vid.author === this.user.id;
				}).map(vid => {
					return Object.assign(
						{
							thumbnailUrl: dummyThumbUrl,
						},
						ContentUtils.expandVideo(Content, vid),
						{
							publishedAt: new Date(vid.publishedAt),
						},
					);
				});

				// Separate the first video from the rest
				this.video = userVideos[0];
				this.otherVideos = userVideos.slice(1);
			},
		},

		watch: {
			'$route.params': function (newParams, oldParams) {
				if (newParams.userName !== oldParams.userName) {
					this.refreshContent();
				}
			},
		},
	}
</script>

<style lang="scss" scoped>
	.container {
		display: flex;
		padding: 0;
		margin: 4rem auto 6rem;
	}

	.banner-section {
		height: 20rem;
		position: relative;
		border-bottom: 0.1rem solid #d6d6d6;

		img {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			object-fit: cover;
		}

		&::after {
			content: '';
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-image: linear-gradient(transparent, 40%, #101010);
			z-index: 1;
		}

		.container {
			margin: 0 auto;
			height: 100%;
			z-index: 100;
			padding: 10rem 0 0;
			display: block;

			h1 {
				color: #eee;
			}
		}
	}

	.video-section {
		flex: 1 1 75%;
	}

	.related-section {
		flex: 1 1 25%;
		padding: 0 0 0 3rem;

		h2 {
			font-size: 2.2rem;
			letter-spacing: -0.08rem;
			line-height: 1.35;
			margin-bottom: 1.5rem;
		}
	}
</style>
