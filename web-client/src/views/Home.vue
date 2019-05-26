<template>
	<div class="home-view container">
		<ViewHeader></ViewHeader>
		<div
			class="category"
			v-for="category in categories">
			<h2>{{ category.title }}</h2>
			<div class="video-list">
				<VideoLink
					:video="video"
					:key="video.title"
					v-for="video in category.videos">
				</VideoLink>
			</div>
		</div>
	</div>
</template>

<script>
	import ViewHeader from '@/components/ViewHeader.vue';
	import VideoLink from '@/components/VideoLink.vue';

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

		mounted () {
			const dummyImgUrl = 'https://dummyimage.com/200x112/000/fff';
			this.categories = [
				{
					title: 'People',
					videos: [
						{
							title: 'Dummy 1',
							thumbnailUrl: dummyImgUrl,
						},
						{
							title: 'Dummy 2',
							thumbnailUrl: dummyImgUrl,
						},
						{
							title: 'Dummy 3',
							thumbnailUrl: dummyImgUrl,
						},
						{
							title: 'Dummy 4',
							thumbnailUrl: dummyImgUrl,
						},
						{
							title: 'Dummy 5',
							thumbnailUrl: dummyImgUrl,
						},
						{
							title: 'Dummy 6',
							thumbnailUrl: dummyImgUrl,
						},
					],
				},
				{
					title: 'Tech',
					videos: [
						{
							title: 'Dummy 1',
							thumbnailUrl: dummyImgUrl,
						},
						{
							title: 'Dummy 2',
							thumbnailUrl: dummyImgUrl,
						},
						{
							title: 'Dummy 3',
							thumbnailUrl: dummyImgUrl,
						},
					],
				},
			];
		},
	}
</script>

<style lang="scss" scoped>
	.category {
		margin: 2rem 0 6rem;
	}

	.video-list {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1fr;
		grid-column-gap: 0.4rem;
		grid-row-gap: 2.4rem;
	}
</style>
