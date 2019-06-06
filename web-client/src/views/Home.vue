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

		async mounted () {
			// Load content (if hasn't already been loaded)
			if (Content.categories.length === 0) {
				await Content.load();
			}
			this.content = Content;
		},

		computed: {
			categories () {
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				if (!this.content) {
					return [];
				}

				return this.content.categories.map(cat => {
					const videos = this.content.videos.filter(vid => {
						return vid.category === cat.id;
					}).map(vid => {
						return Object.assign(
							{
								thumbnailUrl: dummyThumbUrl,
							},
							ContentUtils.expandVideo(Content, vid),
						);
					});
					return Object.assign({}, cat, { videos });
				});
			},

			videoListClasses () {
				return {
					'video-list': true,
					'video-list-multi': this.$breakpoint.name !== 'small',
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

		hr {
			margin: 6rem 0 2rem;
		}

		.video-list {
			display: grid;
			grid-template-columns: 1fr;
			grid-row-gap: 2.4rem;

			&.video-list-multi {
				grid-column-gap: 0.4rem;
				grid-template-columns: 1fr 1fr 1fr 1fr;
			}
		}
	}
</style>
