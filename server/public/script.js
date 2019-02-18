var app = new Vue({
    el: '#app',
    data: {
        sources: null,
        sections: {},
        current: {
            source: null,
            section: null,
            lastItem: null,
            count: 30,
            posts: null,
        },
        isLoading: false,
        error: null
    },
    methods: {},
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


    },
    selectSource: async function(sourceName) {
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

    selectSection: async function(sectionName) {
        if (!this.sections[this.current.source][sectionName]) {
            console.warn(`section ${sectionName} is not present in source ${this.current.source}`);
        }

        this.current.section = sectionName;
        this.getPosts();
    },

    getPosts: async function(loadMore=false) {
        this.isLoading = true;
        this.error = null;

        if (!loadMore) {
            this.current.posts = [];
        }

        try {
            const resp = await axios.get(`/sources/${this.current.source}/sections/${this.current.section}?count=${this.current.count}`);
            this.isLoading = false;
            this.current.posts = resp.data.posts;
        } catch (ex) {
            this.isLoading = false;
            this.error = {
                title: "Failed to get posts",
                message: formatErrorMessage(ex)
            };
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

