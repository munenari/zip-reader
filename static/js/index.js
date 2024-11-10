
window.addEventListener( 'alpine:init', () => {

	async function fetchDirInfo() {
		const hashes = location.hash.substring( 1 ).split( ':' )
		let filepath = hashes[ 0 ] || ''
		const resp = await fetch( `./d/${ filepath }` )
		if (resp.status/100 != 2) {
			throw new Error(`unexpected status code ${resp.status}`)
		}
		return await resp.json()
	}

	Alpine.data('listPage', () => ({
		records: fetchDirInfo(),
		async reload () {
			this.records = await fetchDirInfo()
		}
	}))

} )
