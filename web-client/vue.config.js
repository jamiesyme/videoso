module.exports = {
	devServer: {
		public: 'localhost:3000',
		watchOptions: {
			poll: true,
			ignored: /node_modules/,
		},
	},
};
