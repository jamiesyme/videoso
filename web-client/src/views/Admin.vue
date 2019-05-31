<template>
	<div class="admin-view container">
		<div class="sidebar">
			<nav>
				<router-link
					class="button button-clear"
					to="/admin/categories">
					Categories
				</router-link>
				<router-link
					class="button button-clear"
					to="/admin/users">
					Users
				</router-link>
				<router-link
					class="button button-clear"
					to="/admin/videos">
					Videos
				</router-link>
			</nav>
			<hr>
			<button
				class="button button-clear"
				:disabled="saving || loading"
				@click="saveContent()">
				Save
			</button>
		</div>

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
												v-model="video.publishedAt">
												<template v-slot:label>
													<div class="label">
														{{ prettyDateTime(video.publishedAt) }}
													</div>
												</template>
											</LabelInput>
										</td>
									</tr>
									<tr>
										<td>View count</td>
										<td>
											<LabelInput
												type="number"
												v-model="video.viewCount" />
										</td>
									</tr>
									<tr>
										<td>Category</td>
										<td>
											<LabelInput
												list="category-datalist"
												v-model="video.category" />
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

				<datalist id="category-datalist">
					<option
						:value="category.title"
						v-for="category in categories" />
				</datalist>
				<datalist id="user-datalist">
					<option
						:value="user.name"
						v-for="user in users" />
				</datalist>
			</section>
		</div>
	</div>
</template>

<script>
	import IonAdd from 'vue-ionicons/dist/md-add';
	import IonClose from 'vue-ionicons/dist/md-close';
	import LabelInput from '@/components/LabelInput';
	import Content from '@/content';

	export default {
		props: {
			section: {
				type: String,
				default: 'categories',
			},
		},

		components: {
			IonAdd,
			IonClose,
			LabelInput,
		},

		data () {
			return {
				loading: false,
				saving: false,
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
					 * @property {string} category - category.title
					 * @property {string} author - user.name
					 */
				],
			};
		},

		async mounted () {
			await this.loadContent();
		},

		methods: {
			async loadContent () {
				this.loading = true;
				await Content.load();
				this.categories = Content.categories.slice();
				this.users = Content.users.slice();
				this.videos = Content.videos.map(v => {
					const pub = new Date(Date.parse(v.publishedAt));
					const cat = this.categories.find(c => c.id === v.category);
					const usr = this.users.find(u => u.id === v.author);
					return {
						id:          v.id,
						title:       v.title,
						description: v.description,
						tags:        v.tags.join(' '),
						publishedAt: pub,
						viewCount:   v.viewCount,
						category:    cat.title,
						author:      usr.name,
					};
				});
				this.loading = false;
			},

			async saveContent () {
				this.saving = true;
				Content.categories = this.categories.slice();
				Content.users = this.users.slice();
				Content.videos = this.videos.map(v => {
					const pub = v.publishedAt.toISOString();
					const cat = this.findCategoryByTitle(v.category);
					const usr = this.findUserByName(v.author);
					return {
						id:          v.id,
						title:       v.title,
						description: v.description,
						tags:        v.tags.split(' '),
						publishedAt: pub,
						viewCount:   v.viewCount,
						category:    cat.id,
						author:      usr.id,
					};
				});
				await Content.save();
				this.saving = false;
			},

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

			findCategoryByTitle (title) {
				const s1 = title.toLowerCase().trim();
				return this.categories.find(category => {
					const s2 = category.title.toLowerCase().trim();
					return s1 === s2;
				});
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

			findUserByName (name) {
				const s1 = name.toLowerCase().trim();
				return this.users.find(user => {
					const s2 = user.name.toLowerCase().trim();
					return s1 === s2;
				});
			},

			addVideo () {
				const nextId = this.videos.reduce((accum, next) => {
					return Math.max(next.id + 1, accum);
				}, 1);
				this.videos.push({
					id: nextId,
					title: 'New Video',
					description: '',
					tags: '',
					publishedAt: new Date,
					viewCount: 0,
					category: null,
					author: null,
				});
			},

			removeVideo (videoId) {
				const index = this.videos.findIndex(c => {
					return c.id === videoId;
				});
				this.videos.splice(index, 1);
			},

			prettyDateTime (date) {
				return date.toLocaleString('en-US', {
					month: 'short',
					day: 'numeric',
					year: 'numeric',
					hour: 'numeric',
					minute: '2-digit',
					hour12: true,
				});
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

	.sidebar {
		width: 25%;
		border-right: 0.1rem solid #eee;
		padding-right: 3rem;
		margin-right: 3rem;

		.button {
			display: block;
			width: 100%;
			text-align: left;
			padding: 0 1rem;
		}

		hr {
			border-color: #eee;
		}
	}

	.sub-view {
		flex: 1;
		margin: 0 3rem 0 0;

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
				padding-left: 2.5rem;
				margin-bottom: 0;
				td {
					border: none;
				}
				td:first-of-type {
					width: 20rem;
				}
			}
		}
	}
</style>
