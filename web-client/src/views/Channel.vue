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
		<div class="container" v-if="userVideos.length > 0">
			<section class="video-section">
				<VideoViewer
					:video="primaryVideo"
					v-if="primaryVideo">
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
				content: null,
			};
		},

		async mounted () {
			// Load content (if hasn't already been loaded)
			if (Content.users.length === 0) {
				await Content.load();
			}
			this.content = Content;
		},

		computed: {
			userName () {
				return this.$route.params.userName;
			},

			user () {
				const dummyBannerUrl = 'https://dummyimage.com/1920x200/eee/e6e6e6';

				if (!this.userName) {
					return null;
				}
				if (!this.content) {
					return null;
				}

				const user = this.content.users.find(user => {
					function tr (str) {
						return str.trim().toLowerCase();
					}
					return tr(this.userName) === tr(user.name);
				});
				if (!user) {
					return null;
				}

				if (!user.bannerUrl) {
					user.bannerUrl = dummyBannerUrl;
				}

				return user;
			},

			userVideos () {
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				if (!this.user) {
					return [];
				}
				if (!this.content) {
					return [];
				}

				return this.content.videos.filter(vid => {
					return vid.author === this.user.id;
				}).map(vid => {
					return Object.assign(
						{
							thumbnailUrl: dummyThumbUrl,
						},
						ContentUtils.expandVideo(this.content, vid),
						{
							publishedAt: new Date(vid.publishedAt),
						},
					);
				});
			},

			primaryVideo () {
				return this.userVideos[0];
			},

			otherVideos () {
				return this.userVideos.slice(1);
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
