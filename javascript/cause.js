class Cause {
    constructor(chess) {
        // set link to chess
        this.chess = chess

        // making cause
        this.container = document.createElement("div")

        this.setupCauseElements()
    }

    setupCauseElements() {
        // set class name
        this.container.classList.add("cause")
        this.container.innerHTML = "&nbsp;"
    }

    // show cause
    showCause(cause) {
        this.container.classList.remove("show")
        this.container.innerHTML = cause
        this.container.offsetWidth
        this.container.classList.add("show")
    }

    // remove cause
    hideCause() {
        this.container.classList.remove("show")
    }
}
