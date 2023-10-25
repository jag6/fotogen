const galleryModal = document.getElementById('gallery-modal');
if (galleryModal) {
    //open modal
    const galleryImages = document.querySelectorAll('.gallery-img');
    galleryImages.forEach((image) => {
        image.addEventListener('click', (n) => {
            galleryModal.style.display = 'flex';
            document.querySelector('body').style.overflow = 'hidden';

            //view clicked image
            const index = [...image.parentElement.children].indexOf(image);
            n = index + 1;
            showImages(imageIndex = n);
        });
    });

    //close modal
    document.getElementById('close-modal').addEventListener('click', () => {
        galleryModal.style.display = 'none';
        document.querySelector('body').style.overflow = 'auto';
    });

    //cycle through modal
    const showImages = (n) => {
        let images = document.querySelectorAll('.modal-img');
        if(n > images.length) { imageIndex = 1 }
        if(n < 1) { imageIndex = images.length }
        images.forEach((image) => {
            image.style.display = 'none';
        }); 
        images[imageIndex - 1].style.display = 'flex';
    };

    let imageIndex = 1;
    showImages(imageIndex);

    document.getElementById('next-modal').addEventListener('click', (n) => {
        n = 1;
        showImages(imageIndex += n);
    });
    document.getElementById('prev-modal').addEventListener('click', (n) => {
        n =-1;
        showImages(imageIndex += n)
    });
}