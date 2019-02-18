const app = new Vue({
    el: '#app',
    data: {
        sources: null,
        sections: {},
        current: {
            source: null,
            section: null,
            lastItem: null,
            maxCount: 60,
            posts: null,
            showNext: false,
        },
        isLoading: false,
        error: null
    },
    methods: {
        selectSource: async function(sourceName) {
            if (sourceName === this.current.source) {
                return
            }

            this.isLoading = true;
            this.error = null;
            this.current.source = sourceName;
            if (!this.sections[sourceName]) {
                try {
                    const resp = await axios.get(`/sources/${sourceName}/sections`);
                    this.sections[sourceName] = resp.data.sections;
                } catch (ex) {
                    this.isLoading = false;
                    this.error = {
                        title: "Failed to get sections list",
                        message: formatErrorMessage(ex)
                    };
                    return;
                }
            }

            this.selectSection(this.sections[sourceName][0]);
        },

        onSectionChange() {
            console.log('selection changed');
            this.getPosts();
        },

        selectSection: async function(sectionName) {
            if (!this.sections[this.current.source].includes(sectionName)) {
                console.warn(`section ${sectionName} is not present in source ${this.current.source}`);
            }

            this.current.section = sectionName;
            this.getPosts();
        },

        getPosts: async function(loadMore=false) {
            this.isLoading = true;
            this.error = null;

            let reqUrl = `/sources/${this.current.source}/sections/${this.current.section}?count=${this.current.maxCount}`;

            if (!loadMore) {
                this.current.posts = [];
            } else {
                reqUrl += `&after=${this.current.lastItem}`;
            }

            try {
                const resp = await axios.get(reqUrl);
                this.isLoading = false;
                this.current.posts.push(...resp.data.posts);
                this.updatePostsUI();
            } catch (ex) {
                this.isLoading = false;
                this.error = {
                    title: "Failed to get posts",
                    message: formatErrorMessage(ex)
                };
            }
        },

        updatePostsUI: function () {
            const posts = this.current.posts;
            if (!posts.length) {
                this.current.showNext = false;
                this.current.lastItem = null;
                return;
            }

            this.current.showNext = (posts.length >= this.current.maxCount);
            this.current.lastItem = posts[posts.length - 1].id;
        }
    },
    created: async function() {
        try {
            this.error = null;
            const resp = await axios.get('/sources');
            this.sources = resp.data.sources;

            if (!this.sources.length) {
                this.isLoading = false;
                return;
            }

            this.selectSource(this.sources[0]);
        } catch (ex) {
            this.isLoading = false;
            this.error = {
                title: "Failed to get sources list",
                message: formatErrorMessage(ex)
            };
            return;
        }


    }
});

const formatErrorMessage = (err) => {
  if (!err.response || !err.response.data) {
      return err.message
  }

  if (err.response.data.error) {
      return `${err.response.data.error} (code: ${err.response.data.code})`
  }

  return err.message
};

