
function getThumbImg ( record, page ) {
	const img = document.createElement( 'img' )
	img.src = `./c/${ record.hashed_name }?page=${ page }`
	img.loading = 'lazy'
	return img
}

function getLinkEl ( record ) {
	const el = document.createElement( 'a' )
	el.className = 'record'
	const imgWrapEl = document.createElement( 'div' )
	el.href = `viewer.html#${ record.hashed_name }`
	if ( record.is_dir ) {
		el.href = `./#${ record.hashed_name }`
	}
	const pEl = document.createElement( 'p' )
	pEl.innerText = record.name
	imgWrapEl.appendChild( getThumbImg( record, 0 ) )
	imgWrapEl.appendChild( getThumbImg( record, 1 ) )
	el.appendChild( imgWrapEl )
	el.appendChild( pEl )
	return el
}

async function getDirInfo () {
	const hashes = location.hash.substring( 1 ).split( ':' )
	let filepath = hashes[ 0 ] || ''
	const resp = await fetch( `./d/${ filepath }` )
	const respData = await resp.json()
	const contentEl = document.getElementsByClassName( 'contents' )[ 0 ]
	contentEl.innerHTML = ''
	const fragment = document.createDocumentFragment()
	respData.forEach( r => {
		const el = getLinkEl( r )
		fragment.appendChild( el )
	} )
	contentEl.appendChild( fragment )
}

window.onhashchange = () => getDirInfo()
getDirInfo()
