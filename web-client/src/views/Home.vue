<template>
	<div class="home-view">
		<div class="container">
			<div
				class="category"
				v-for="category in categories">
				<h2>{{ category.title }}</h2>
				<div :class="videoListClasses">
					<VideoLink
						:video="video"
						:key="video.id"
						:extraWidth="videoLinkExtraWidth"
						v-for="video in category.videos">
					</VideoLink>
				</div>
				<hr>
			</div>
		</div>
	</div>
</template>

<script>
	import Content from '@/content';
	import ContentUtils from '@/utils/content';
	import ViewHeader from '@/components/ViewHeader';
	import VideoLink from '@/components/VideoLink';

	function expandVideo (video) {
		return ContentUtils.expandVideo(Content, video);
	}

	export default {
		name: 'home',

		components: {
			ViewHeader,
			VideoLink,
		},

		data () {
			return {
				content: null,
			};
		},

		computed: {
			categories () {
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				return Content.categories.map(cat => {
					const videos = Content.videos.filter(vid => {
						return vid.category === cat.id;
					}).map(vid => {
						return Object.assign(
							{
								thumbnailUrl: dummyThumbUrl,
							},
							expandVideo(vid),
						);
					});
					return Object.assign({}, cat, { videos });
				});
			},

			videoListClasses () {
				const bp = this.$breakpoint.name;
				const multi2 = bp === 'medium';
				const multi4 = bp === 'large' || bp === 'xlarge';
				return {
					'video-list': true,
					'video-list-multi-2': multi2,
					'video-list-multi-4': multi4,
				};
			},

			videoLinkExtraWidth () {
				if (this.$breakpoint.name === 'small') {
					return '4rem';
				} else {
					return null;
				}
			},
		},
	}
</script>

<style lang="scss" scoped>
	.category {
		margin: 2rem 0;

		&:first-of-type {
			margin-top: 4rem;
		}
		&:last-of-type {
			margin-bottom: 6rem;
			hr {
				display: none;
			}
		}

		h2 {
			margin: 0 0 3rem;
		}

		hr {
			margin: 6rem 0 2rem;
		}

		.video-list {
			display: grid;
			grid-template-columns: 1fr;
			grid-row-gap: 4rem;

			&.video-list-multi-2 {
				grid-column-gap: 0.4rem;
				grid-row-gap: 2.4rem;
				grid-template-columns: 1fr 1fr;
			}
			&.video-list-multi-4 {
				grid-column-gap: 0.4rem;
				grid-row-gap: 2.4rem;
				grid-template-columns: 1fr 1fr 1fr 1fr;
			}
		}
	}
</style>
