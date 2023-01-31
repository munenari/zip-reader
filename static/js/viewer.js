
const mainTagName = 'img'
const mainTagProp = 'src'

let filepath
let page
let prevPage = -1
let maxPage = 9999
let displayPage = 2
let prevFilePath
let nextFilePath

function setPage ( p ) {
	if ( p < 0 ) {
		p = 0
	}
	if ( p >= maxPage ) {
		p = maxPage - 1
	}
	localStorage.setItem( 'page:' + filepath, p )
	page = p
	setPageFooter( page )
	setImgSrc( page )
}

function setImgTags ( num ) {
	localStorage.setItem( 'imgtags', num )
	const createEl = ( el, n ) => {
		el.innerHTML = ''
		for ( let i = 0; i < n; i++ ) {
			const child = document.createElement( mainTagName )
			child.className = 'img'
			el.appendChild( child )
		}
	}
	createEl( document.getElementsByClassName( 'contents' )[ 0 ], num )
	createEl( document.getElementsByClassName( 'preload' )[ 0 ], 2 )
	setPage( page )
}

function clock () {
	const date = new Date()
	const dateH = ( '0' + date.getHours() ).slice( -2 )
	const dateM = ( '0' + date.getMinutes() ).slice( -2 )
	const el = document.querySelector( '.header .clock' )
	el.textContent = `${ dateH }:${ dateM }`
}

function pageSize ( el ) {
	return new Promise( resolve => {
		el.onload = () => {
			el.onload = undefined
			resolve( { width: el.width, height: el.height } )
		}
	} )
}

function getContentURL ( p ) {
	return `./c/${ filepath }?page=${ p }`
}

function setImgSrc ( p ) {
	const els = document.querySelectorAll( '.contents .img' )
	const fetchPromises = []
	els.forEach( ( el, index ) => {
		fetchPromises.push( pageSize( el ) )
		el[ mainTagProp ] = getContentURL( p + index )
		if ( p + index >= maxPage ) {
			el[ mainTagProp ] = ''
		}
		if ( p + index < 0 ) {
			el[ mainTagProp ] = ''
		}
	} )
	const hideEls = document.querySelectorAll( '.preload .img' )
	hideEls.forEach( ( el, index ) => {
		el[ mainTagProp ] = getContentURL( p + index + els.length )
		if ( p + index >= maxPage ) {
			el[ mainTagProp ] = ''
		}
		if ( p + index < 0 ) {
			el[ mainTagProp ] = ''
		}
	} )
	const hideImgPos = prevPage <= p ? 1 : 0
	Promise.all( fetchPromises ).then( values => {
		displayPage = 2
		els.forEach( el => {
			el.hidden = false
		} )
		values.forEach( v => {
			if ( v.width > v.height ) {
				displayPage = 1
				els[ hideImgPos ].hidden = true
			}
		} )
	} )
	prevPage = p
}

async function getFileInfo () {
	const h = location.hash.substring( 1 )
	if ( !h ) {
		return
	}
	setImgSrc( -999 )
	setHeader( { name: '/Loading...' } )
	setPageFooter( -1 )
	filepath = h
	page = parseInt( localStorage.getItem( 'page:' + filepath ) ) || 0
	const resp = await fetch( `./i/${ filepath }` )
	const respData = await resp.json()
	maxPage = respData.size
	prevFilePath = respData.prev_hashed_name
	nextFilePath = respData.next_hashed_name
	console.log( respData )
	setImgSrc( page )
	setHeader( respData )
	setPageFooter( page )
}

function setHeader ( respData ) {
	const el = document.querySelector( '.header .title' )
	el.textContent = respData.name.split( '/' ).pop().replace( /\.zip$/, '' )
	const aEl = document.querySelector( '.header a' )
	const prev = aEl.href
	aEl.href = `./#${ respData.parent_hashed_name }`
	if ( !respData.parent_hashed_name ) {
		aEl.href = prev
	}
}

function setPageFooter ( page ) {
	const el = document.getElementsByClassName( 'page' )[ 0 ]
	el.innerHTML = ''
	el.textContent = `[${ page + 1 }/${ maxPage }]`
	if ( page < 0 ) {
		el.textContent = '[-/-]'
	}
}

function checkKey ( e ) {
	e = e || window.event

	if ( e.shiftKey || e.metaKey || e.ctrlKey || e.altKey ) {
		return
	}
	if ( e.keyCode == '38' ) { // up
		location.hash = prevFilePath
		getFileInfo()
	}
	else if ( e.keyCode == '40' ) { // down
		location.hash = nextFilePath
		getFileInfo()
	}
	else if ( e.keyCode == '37' || e.keyCode == '32' ) { // left, space
		setPage( page + displayPage )
	}
	else if ( e.keyCode == '39' ) { // right
		setPage( page - displayPage )
	}
	else if ( e.keyCode == '71' ) { // G
		setPage( page - 1 )
	}
	else if ( e.keyCode == '70' ) { // F
		toggleFullscreen()
	}
	else if ( e.keyCode == '49' ) { // 1
		setImgTags( 1 )
	}
	else if ( e.keyCode == '50' ) { // 2
		setImgTags( 2 )
	}
	console.log( e.keyCode, e )
}

function toggleFullscreen () {
	if ( !document.fullscreenElement ) {
		const view = document.getElementsByClassName( 'main' )[ 0 ]
		view.requestFullscreen()
	} else {
		document.exitFullscreen()
	}
}

window.onhashchange = () => getFileInfo()
document.onkeydown = checkKey
getFileInfo()
clock()
setImgTags( localStorage.getItem( 'imgtags' ) || 1 )
setInterval( () => clock(), 10 * 1000 )

let timer
window.addEventListener( 'mousemove', () => {
	document.body.classList.remove( "hide-cursor" )
	clearTimeout( timer )
	timer = setTimeout( () => document.body.classList.add( "hide-cursor" ), 3000 )
} )
