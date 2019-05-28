<template>
	<div id="app">
		<ViewHeader></ViewHeader>
		<router-view></router-view>
	</div>
</template>

<script>
	import ViewHeader from '@/components/ViewHeader';

	export default {
		components: {
			ViewHeader,
		},
	}
</script>

<style lang="scss">
	@import url('//fonts.googleapis.com/css?family=Roboto:300,300italic,700,700italic');
	@import url('//cdn.rawgit.com/necolas/normalize.css/master/normalize.css');
	@import url('//cdn.rawgit.com/milligram/milligram/master/dist/milligram.min.css');
	@import '~vue-ionicons/ionicons.scss';

	#app {
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
		color: #2c3e50;
	}
	#nav {
		padding: 30px;
		a {
			font-weight: bold;
			color: #2c3e50;
			&.router-link-exact-active {
				color: #42b983;
			}
		}
	}
</style>
