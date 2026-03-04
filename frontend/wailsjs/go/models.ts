export namespace sysinfo {

	export class Info {
	    cpu: string;
	    memory: string;
	    gpu: string;
	    disk: string;

	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.memory = source["memory"];
	        this.gpu = source["gpu"];
	        this.disk = source["disk"];
	    }
	}

}
