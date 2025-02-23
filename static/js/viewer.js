
window.addEventListener('alpine:init', () => {

	async function getFileInfoData (h) {
		if ( !h ) {
			return
		}
		const resp = await fetch( `./i/${ h }` )
		return await resp.json()
	}

	Alpine.data('clock', () => ({
		now: new Date(),
		_timer: null,
		init () {
			this._timer = setInterval( () => {
				this.now = new Date()
			}, 10 * 1000 )
		},
		destroy () {
			clearInterval(this._timer)
		},
		clockText () {
			const dateH = ( '0' + this.now.getHours() ).slice( -2 )
			const dateM = ( '0' + this.now.getMinutes() ).slice( -2 )
			return `${ dateH }:${ dateM }`
		}
	}))

	let DEFAULT_DISPLAYSIZE = 1
	if (window.innerWidth > 768) {
		DEFAULT_DISPLAYSIZE = 2
	}
	Alpine.data('viewPage', () => ({
		filePath: '',
		fileInfo: {
			name: '',
			parent_hashed_name: '',
			size: -1
		},
		page: 0,
		displayPage: parseInt( localStorage.getItem( 'displaySize' ) ) || DEFAULT_DISPLAYSIZE,
		oddNum: 0,
		mouseMoving: false,
		mouseMoveTimer: null,
		showPageSlider: false,

		title () {
			const name = this.fileInfo.name || 'Loading...'
			return name.split( '/' ).pop().replace( /\.zip$/, '' )
		},
		pageText () {
			const p = this.page * this.displayPage
			const s = this.fileInfo.size
			return `[${p}/${s}]`
		},
		async init () {
			this.filePath = location.hash.substring( 1 )
			this.fileInfo = await getFileInfoData(this.filePath)
			this.$watch('page', (v) => {
				const el = document.getElementById(`page${v}`)
				if (!el) return
				el.scrollIntoView( { behavior: 'auto' } )
				localStorage.setItem( 'page:' + this.filePath, v )
			})
			this.page = parseInt( localStorage.getItem( 'page:' + this.filePath ) ) || 0
			console.info(this.$data.fileInfo)
		},
		getContentURL ( p ) {
			const eagerSize = 2
			if (Math.abs(p/this.displayPage-this.page) > eagerSize) {
				return ''
			}
			return `./c/${ this.filePath }?page=${ p }`
		},
		setPage (v) {
			this.page = v
			if (this.page > this.fileInfo.size / this.displayPage) {
				this.page = this.fileInfo.size / this.displayPage
			}
			if (this.page < 0) {
				this.page = 0
			}
		},
		prevPage () {
			this.setPage(this.page-1)
		},
		nextPage () {
			this.setPage(this.page+1)
		},
		async toggleFullscreen () {
			const current = this.page
			if ( !document.fullscreenElement ) {
				await document.body.requestFullscreen()
			} else {
				await document.exitFullscreen()
			}
			this.page = current
		},
		setDisplayPage (v) {
			this.displayPage = v
			localStorage.setItem( 'displaySize', v )
		},
		changeHash (v) {
			if (!v) return
			location.hash = v
		},
		oddPage () {
			if (this.oddNum == 0) {
				this.oddNum = -1
			} else {
				this.oddNum = 0
			}
		},
		onMousemove () {
			this.mouseMoving = true
			clearTimeout(this.mouseMoveTimer)
			this.mouseMoveTimer = setTimeout( () => {
				this.mouseMoving = false
			}, 3000 )
		},
		onInputPage (evt) {
			this.page = -evt.target.value
		},
		getImgLoading (p) {
			const eagerSize = 2
			if (Math.abs(p-this.page) <= eagerSize) {
				return 'eager'
			}
			return 'lazy'
		},
		deferSetPage(i) {
			clearTimeout(this.setPageTimer)
			this.setPageTimer = setTimeout(() => this.setPage(i), 200)
		}
	}))

})
