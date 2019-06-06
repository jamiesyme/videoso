<template>
	<div class="search-view">
		<div class="container">
			<h1 v-if="videos.length > 0">
				Showing search results for: {{ query }}
			</h1>
			<h1 v-else>
				No search results for: {{ query }}
			</h1>
			<div :class="videoListClasses">
				<VideoLink
					:video="video"
					:key="video.id"
					:extraWidth="videoLinkExtraWidth"
					v-for="video in videos">
				</VideoLink>
			</div>
		</div>
	</div>
</template>

<script>
	import FlexSearch from 'flexsearch';
	import Content from '@/content';
	import ContentUtils from '@/utils/content';
	import VideoLink from '@/components/VideoLink';

	export default {
		components: {
			VideoLink,
		},

		data () {
			return {
				content: null,
			};
		},

		async mounted () {
			// Load content (if hasn't already been loaded)
			if (Content.videos.length === 0) {
				await Content.load();
			}
			this.content = Content;
		},

		computed: {
			query () {
				return this.$route.query.q;
			},

			searchIndex () {
				if (!this.content) {
					return null;
				}

				const searchIndex = new FlexSearch('match', {
					doc: {
						id: 'id',
						field: [
							'title',
							'description',
							'tags',
							'category',
							'author',
						],
					},
				});
				searchIndex.add(this.content.videos.map(vid => {
					const fullVid = ContentUtils.expandVideo(this.content, vid);
					return {
						id:          fullVid.id,
						title:       fullVid.title,
						description: fullVid.description,
						tags:        fullVid.tags.join(' '),
						category:    fullVid.category.title,
						author:      fullVid.author.name,
					};
				}));
				return searchIndex;
			},

			videos () {
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				// We need an index (and content) to search
				if (!this.searchIndex) {
					return [];
				}
				if (!this.content) {
					return [];
				}

				// Search videos
				const results = this.searchIndex.search(this.query);

				// Translate search results
				return results.map(searchVid => {
					const vid = this.content.videos.find(v => v.id === searchVid.id);
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
	.search-view {
		margin: 4rem 0 6rem;
	}

	h1 {
		font-size: 2.4rem;
		line-height: 1.3;
		margin: 3rem 0;
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
</style>
