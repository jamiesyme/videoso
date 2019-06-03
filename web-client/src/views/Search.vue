<template>
	<div class="search-view">
		<div class="container">
			<h1 v-if="videos.length > 0">
				Showing search results for: {{ query }}
			</h1>
			<h1 v-else>
				No search results for: {{ query }}
			</h1>
			<div class="video-list">
				<VideoLink
					:video="video"
					:key="video.id"
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
				searchIndex: null,
			};
		},

		mounted () {
			this.refreshContent();
		},

		methods: {
			async refreshContent () {

				// Load content (if hasn't already been loaded)
				if (Content.videos.length === 0) {
					await Content.load();
				}

				// Build the search index
				this.searchIndex = new FlexSearch('match', {
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
				this.searchIndex.add(Content.videos.map(vid => {
					const fullVid = ContentUtils.expandVideo(Content, vid);
					return {
						id:          fullVid.id,
						title:       fullVid.title,
						description: fullVid.description,
						tags:        fullVid.tags.join(' '),
						category:    fullVid.category.title,
						author:      fullVid.author.name,
					};
				}));
			},
		},

		computed: {
			query () {
				return this.$route.query.q;
			},

			videos () {
				const dummyThumbUrl = 'https://dummyimage.com/200x112/000/fff';

				// We need an index to search
				if (!this.searchIndex) {
					return [];
				}

				// Search videos
				const results = this.searchIndex.search(this.query);

				// Translate search results
				return results.map(searchVid => {
					const vid = Content.videos.find(v => v.id === searchVid.id);
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
			},
		},
	}
</script>

<style lang="scss" scoped>
	.container {
		padding: 0;
		margin: 4rem auto 6rem;
	}

	h1 {
		font-size: 2.4rem;
		line-height: 1.3;
		margin: 3rem 0;
	}

	.video-list {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-column-gap: 0.4rem;
		grid-row-gap: 2.4rem;
	}
</style>
