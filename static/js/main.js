//MOBILE NAV
const mobile_nav = document.getElementById('mobile-nav');
const nav_overlay = document.getElementById('nav-overlay');
const body = document.querySelector('body');

const clickOff = () => {
    nav_overlay.style.display = 'none';
    mobile_nav.style.width = '0%';
    body.style.overflow = 'auto';
};
document.getElementById('hamburger-icon').addEventListener('click', () => {
    if(mobile_nav.style.width === '300px') {
        clickOff();
    }else {
        mobile_nav.style.width = '300px';
        nav_overlay.style.display = 'flex';
        body.style.overflow = 'hidden';
    }
});
nav_overlay.addEventListener('click', (e) => {
    if(e.target == nav_overlay) {
        clickOff();
    }
});


//ANIMATIONS
const scrollElements = document.querySelectorAll('.js-scroll');
var throttleTimer;
const throttle = (callback, time) => {
    if (throttleTimer) return;
    throttleTimer = true;
    setTimeout(() => {
        callback();
        throttleTimer = false;
        }, time
    );
}
const elementInView = (el, dividend = 1) => {
const elementTop = el.getBoundingClientRect().top;
	return (
    	elementTop <= (window.innerHeight || document.documentElement.clientHeight) / dividend
	);
};
const elementOutofView = (el) => {
	const elementTop = el.getBoundingClientRect().top;
	return (
    	elementTop > (window.innerHeight || document.documentElement.clientHeight)
  );
};
const displayScrollElement = (element) => {
  	element.classList.add('scrolled');
};
const hideScrollElement = (element) => {
  	element.classList.remove('scrolled');
};
const handleScrollAnimation = () => {
  	scrollElements.forEach((el) => {
    	if (elementInView(el, 1.25)) {
      	displayScrollElement(el);
    	}else if (elementOutofView(el)) {
      	hideScrollElement(el)
    }
  })
};
window.addEventListener("scroll", () => { 
    throttle(() => {
        handleScrollAnimation();
    }, 250);
});


//ALERTS
document.querySelectorAll('.alert-message').forEach((alert) => {
	//auto delete
	setTimeout(() => {
		alert.remove();
	}, 5000);
	//click to delete
	alert.addEventListener('click', () => {
		alert.classList.add('fade-out');
		setTimeout(() => {
			alert.remove();
		}, 300);
	});
});


//DELETE GALLERY
if(document.querySelector('.delete-gallery-btn')) {
	document.querySelectorAll('.delete-gallery-btn').forEach((btn) => {
		btn.addEventListener('click', (e) => {
			const title = e.target.getAttribute('title');
			if(title) {
				if(!confirm('Do you want to delete your gallery: ' + title + '?')) {
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

//SHOW GALLERY 
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


//HOME
const changeHomeText = document.getElementById('home-header');
const splitScreenBackground = document.getElementById('split-screen');

const shuffleColors = (array) => {
    let currentIndex = array.length;

    while(currentIndex != 0){
        let randomIndex = Math.floor(Math.random() * array.length);
        currentIndex -=1;

        let temp = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temp;
    }

    return array;
};
const setSplitScreenBackgroundColor = () => {
    let colors = [
        '#1A88BD', '#0A4F70', '#BDB71A', '#3D0209', '#3D151A',
        '#C84453', '#007018', '#FF243D', '#BD961A', '#573326', '#8A0615'
    ]; 
    shuffleColors(colors);
    splitScreenBackground.style.backgroundImage = `linear-gradient(to right, ${colors[0]}, ${colors[1]})`; 
    // splitScreenBackground.style.backgroundColor = colors[0]; 
}

if(changeHomeText) {
    setTimeout(() => {
        changeHomeText.querySelector('#white-text').textContent = 'YOUR';
    }, 5000);
    setTimeout(() => {
        changeHomeText.querySelector('#white-text').textContent = 'OUR';
    }, 6500);
    setTimeout(() => {
        setInterval(() => {
            setSplitScreenBackgroundColor();
        }, 1500);
    }, 6500);
}

//RENDER
// document.querySelector('.render').innerHTML = `
   
// `;