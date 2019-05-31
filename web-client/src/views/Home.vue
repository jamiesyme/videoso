<template>
	<div class="home-view">
		<div class="container">
			<div
				class="category"
				v-for="category in categories">
				<h2>{{ category.title }}</h2>
				<div class="video-list">
					<VideoLink
						:video="video"
						:key="video.id"
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
				categories: [],
			};
		},

		async mounted () {
			const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

			// Load content (if hasn't already been loaded)
			if (Content.categories.length < 1) {
				await Content.load();
			}

			// Fill categories from content
			this.categories = Content.categories.map(cat => {
				return Object.assign({}, cat, {
					videos: Content.videos.filter(vid => {
						return vid.category === cat.id;
					}).map(vid => {
						return Object.assign(
							ContentUtils.expandVideo(Content, vid),
							{
								thumbnailUrl: dummyThumbUrl,
							},
						);
					}),
				});
			});
		},
	}
</script>

<style lang="scss" scoped>
	.container {
		padding: 0;
	}

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
			grid-template-columns: 1fr 1fr 1fr 1fr;
			grid-column-gap: 0.4rem;
			grid-row-gap: 2.4rem;
		}

		.video-link {
			font-size: 1.4rem;
		}
	}
</style>
