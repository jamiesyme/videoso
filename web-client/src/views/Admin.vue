<template>
	<div class="admin-view container">
		<nav>
			<button
				class="button-clear"
				@click="section = 'categories'">
				Categories
			</button>
			<button
				class="button-clear"
				@click="section = 'users'">
				Users
			</button>
			<button
				class="button-clear"
				@click="section = 'videos'">
				Videos
			</button>
		</nav>
		<div class="sub-view">
			<section
				class="categories-section"
				v-if="section === 'categories'">
				<h2>Categories</h2>
				<table>
					<thead>
						<tr>
							<th class="title-col">Title</th>
							<th class="delete-col">Delete</th>
						</tr>
					</thead>
					<tbody>
						<tr
							class="category"
							v-for="category in categories">
							<td>
								<LabelInput v-model="category.title">
								</LabelInput>
							</td>
							<td>
								<button
									class="button-clear"
									@click="removeCategory(category.id)">
									<IonClose />
								</button>
							</td>
						</tr>
						<tr class="new-category">
							<td>
								<button
									class="button-clear"
									@click="addCategory()">
									<IonAdd />Add category
								</button>
							</td>
							<td></td>
						</tr>
					</tbody>
				</table>
			</section>
			<section
				class="users-section"
				v-if="section === 'users'">
				<h2>Users</h2>
			</section>
			<section
				class="videos-section"
				v-if="section === 'videos'">
				<h2>Videos</h2>
			</section>
		</div>
	</div>
</template>

<script>
	import IonAdd from 'vue-ionicons/dist/md-add';
	import IonClose from 'vue-ionicons/dist/md-close';
	import LabelInput from '@/components/LabelInput';

	export default {
		components: {
			IonAdd,
			IonClose,
			LabelInput,
		},

		data () {
			return {
				section: 'categories',
				categories: [
					{
						id: 1,
						title: 'Tech',
					},
					{
						id: 2,
						title: 'People',
					},
				],
			};
		},

		methods: {
			addCategory () {
				const nextId = this.categories.reduce((accum, next) => {
					return Math.max(next.id + 1, accum);
				}, 1);
				this.categories.push({
					id: nextId,
					title: 'New Category',
				});
			},

			removeCategory (categoryId) {
				const index = this.categories.findIndex(c => {
					return c.id === categoryId;
				});
				this.categories.splice(index, 1);
			},
		},
	}
</script>

<style lang="scss" scoped>
	.admin-view {
		display: flex;
		margin: 4rem auto 6rem;
		padding: 0;
	}

	nav {
		width: 25%;
		border-right: 0.1rem solid #eee;

		button {
			display: block;
			width: 100%;
			text-align: left;
			padding: 0 1rem;
		}
	}

	.sub-view {
		flex: 1;
		margin: 0 0 0 3rem;

		table {
			.title-col {
				width: 100%;
			}

			tr {
				button {
					padding: 0;
					margin: 0;
				}
				.ion {
					font-size: 1.6rem;
				}
				&.new-category {
					.ion {
						position: relative;
						top: 0.3rem;
						margin-right: 0.2rem;
					}
				}
				&.category {
					button {
						width: 100%;
					}
					.ion {
						display: block;
					}
				}
			}
		}
	}
</style>
