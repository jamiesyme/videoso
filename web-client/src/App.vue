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
	@import '~vue-ionicons/ionicons.scss';

	$color-primary: #3bacd2;
	@import '~milligram/src/milligram.sass';

	$plyr-color-main: #40bae3;
	@import '~plyr/src/sass/plyr.scss';

	#app {
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
	}

	hr {
		border-color: #eee;
	}

	.aspect-16-9 {
		width: 100%;
		&::before {
			content: '';
			width: 100%;
			padding-top: 56.25%;
			float: left;
		}
		&::after {
			content: '';
			clear: left;
			display: block;
		}
	}

	.plyr {

		button {
			margin-bottom: 0;
			height: auto;
			text-transform: none;

			&:focus {
				background-color: inherit;
				border-color: inherit;
			}
		}

		.plyr__menu__container [role=menu] {
			margin-top: -2px;
		}
	}
</style>
