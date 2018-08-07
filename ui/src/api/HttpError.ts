export class HttpError extends Error {
    public readonly status: number;

    constructor(status: number, statusText: string) {
        super(statusText);
        this.status = status;

        // Set the prototype explicitly.
        Object.setPrototypeOf(this, HttpError.prototype);
    }
}
