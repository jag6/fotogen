if(document.querySelector('.delete-gallery-btn')) {
	document.querySelectorAll('.delete-gallery-btn').forEach((btn) => {
		btn.addEventListener('click', (e) => {
			const title = e.target.getAttribute('title');
			if(title) {
				if(!confirm(`Do you want to delete your gallery: "${title}"?`)) {
					e.preventDefault();
				}
			}
			else {
				if(!confirm('Do you want to delete this image?')) {
					e.preventDefault();
				}
			}
		});
	});
}