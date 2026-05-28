
window.addEventListener( 'alpine:init', () => {

	async function fetchDirInfo () {
		const paths = [ './d' ]
		const hashes = location.hash.substring( 1 ).split( ':' )
		if ( hashes[ 0 ] ) {
			paths.push( hashes[ 0 ] )
		}
		const resp = await fetch( `${ paths.join( '/' ) }` )
		if ( resp.status / 100 != 2 ) {
			throw new Error( `unexpected status code ${ resp.status }` )
		}
		return await resp.json()
	}

	Alpine.data( 'listPage', () => ( {
		records: fetchDirInfo(),
		async reload () {
			location.reload()
		}
	} ) )

} )
