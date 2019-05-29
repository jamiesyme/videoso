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
								<LabelInput v-model="category.title" />
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
									<IonAdd />Add Category
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
				<table>
					<thead>
						<tr>
							<th class="name-col">Name</th>
							<th class="banner-col">Banner URL</th>
							<th class="delete-col">Delete</th>
						</tr>
					</thead>
					<tbody>
						<tr
							class="user"
							v-for="user in users">
							<td>
								<LabelInput v-model="user.name" />
							</td>
							<td class="banner-td">
								<LabelInput v-model="user.bannerUrl">
									<template v-slot:empty>
										<span class="null-banner">none</span>
									</template>
								</LabelInput>
							</td>
							<td>
								<button
									class="button-clear"
									@click="removeUser(user.id)">
									<IonClose />
								</button>
							</td>
						</tr>
						<tr class="new-user">
							<td>
								<button
									class="button-clear"
									@click="addUser()">
									<IonAdd />Add User
								</button>
							</td>
							<td></td>
						</tr>
					</tbody>
				</table>
			</section>

			<section
				class="videos-section"
				v-if="section === 'videos'">
				<h2>Videos</h2>
				<table>
					<thead>
						<tr>
							<th class="video-col">Video</th>
							<th class="delete-col">Delete</th>
						</tr>
					</thead>
					<tbody>
						<tr
							class="video"
							v-for="video in videos">
							<td>
								<table>
									<tr>
										<td>Title</td>
										<td>
											<LabelInput v-model="video.title" />
										</td>
									</tr>
									<tr>
										<td>Description</td>
										<td>
											<LabelInput v-model="video.description" />
										</td>
									</tr>
									<tr>
										<td>Tags</td>
										<td>
											<LabelInput v-model="video.tags" />
										</td>
									</tr>
									<tr>
										<td>Published at</td>
										<td>
											<LabelInput
												type="datetime-local"
												v-model="video.publishedAt" />
										</td>
									</tr>
									<tr>
										<td>View count</td>
										<td>
											<LabelInput
												type="number"
												v-model="video.views" />
										</td>
									</tr>
									<tr>
										<td>Author</td>
										<td>
											<LabelInput
												list="user-datalist"
												v-model="video.author" />
										</td>
									</tr>
								</table>
							</td>
							<td class="delete-td">
								<button
									class="button-clear"
									@click="removeVideo(video.id)">
									<IonClose />
								</button>
							</td>
						</tr>
						<tr class="new-video">
							<td>
								<button
									class="button-clear"
									@click="addVideo()">
									<IonAdd />Add Video
								</button>
							</td>
							<td></td>
						</tr>
					</tbody>
				</table>

				<datalist id="user-datalist">
					<option
						:value="user.name"
						v-for="user of users" />
				</datalist>
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
					/**
					 * @typedef {object} Category
					 * @property {number} id
					 * @property {string} title
					 */
				],
				users: [
					/**
					 * @typedef {object} User
					 * @property {number} id
					 * @property {string} name
					 * @property {string} [bannerUrl]
					 */
				],
				videos: [
					/**
					 * @typedef {object} Video
					 * @property {number} id
					 * @property {string} title
					 * @property {string} description
					 * @property {string} tags - space-separated
					 * @property {string} publishedAt - YYYY-MM-DDThh:mm
					 * @property {number} viewCount
					 * @property {string} author - user.name
					 */
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

			addUser () {
				const nextId = this.users.reduce((accum, next) => {
					return Math.max(next.id + 1, accum);
				}, 1);
				this.users.push({
					id: nextId,
					name: 'New User',
					bannerUrl: null,
				});
			},

			removeUser (userId) {
				const index = this.users.findIndex(c => {
					return c.id === userId;
				});
				this.users.splice(index, 1);
			},

			addVideo () {
				function getTodayDateTimeLocal () {
					const date = new Date;
					const offset = date.getTimezoneOffset();
					const msPerMin = 60 * 1000;
					date.setTime(date.getTime() - offset * msPerMin);
					return date.toISOString().substr(0, 16);
				}
				const nextId = this.videos.reduce((accum, next) => {
					return Math.max(next.id + 1, accum);
				}, 1);
				const todayDate = getTodayDateTimeLocal();
				this.videos.push({
					id: nextId,
					title: 'New Video',
					description: '',
					tags: '',
					publishedAt: todayDate,
					viewCount: 0,
					author: null,
				});
			},

			removeVideo (videoId) {
				const index = this.videos.findIndex(c => {
					return c.id === videoId;
				});
				this.videos.splice(index, 1);
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

		.categories-section, .users-section, .videos-section {
			table {
				button {
					padding: 0;
					margin: 0;
				}
				.ion {
					font-size: 2.2rem;
				}
				.category, .user, .video {
					button {
						width: 100%;
					}
					.ion {
						display: block;
					}
				}
				.new-category, .new-user, .new-video {
					.ion {
						position: relative;
						top: 0.6rem;
						margin: 0 0.6rem;
					}
				}
			}
		}

		.categories-section {
			.title-col {
				width: 100%;
			}
		}

		.users-section {
			.name-col {
				width: 50%;
			}
			.banner-col {
				width: 50%;
			}
			.banner-td {
				max-width: 50rem;
			}
			.null-banner {
				font-style: italic;
			}
		}

		.videos-section {
			.video-col {
				width: 100%;
			}
			.delete-td {
				vertical-align: top;
			}
			.video table {
				margin-bottom: 0;
				td {
					border: none;
				}
			}
		}
	}
</style>
