//Animations
const scrollElements = document.querySelectorAll(".js-scroll");
const elementInView = (el, dividend = 1) => {
  const elementTop = el.getBoundingClientRect().top;
  	return (
    	elementTop <=
    	(window.innerHeight || document.documentElement.clientHeight) / dividend
	);
};
const elementOutofView = (el) => {
  	const elementTop = el.getBoundingClientRect().top;
  	return (
    	elementTop > (window.innerHeight || document.documentElement.clientHeight)
  	);
};
const displayScrollElement = (element) => {
  	element.classList.add("scrolled");
};
const hideScrollElement = (element) => {
 	element.classList.remove("scrolled");
};
const handleScrollAnimation = () => {
	scrollElements.forEach((el) => {
    	if (elementInView(el, 1.25)) {
      	displayScrollElement(el);
    	} else if (elementOutofView(el)) {
      	hideScrollElement(el)
    	}
  	});
};
window.addEventListener("scroll", () => { 
	handleScrollAnimation();
});


//Color Switch
const isLight = () => {
    return localStorage.getItem('light-mode');
};
const toggleRootClass = () => {
    document.querySelector(':root').classList.toggle('light');
};
const toggleLocalStorageItem = () => {
    if (isLight()) {
        localStorage.removeItem('light-mode');
    }else {
        localStorage.setItem('light-mode', 'set');
    }
};
if(isLight()) {
    toggleRootClass();
}
document.querySelector('.theme-icon').addEventListener('click', () => {
    toggleLocalStorageItem();
    toggleRootClass();
});