export namespace runner {

	export class AppItem {
	    name: string;
	    id: string;

	    static createFrom(source: any = {}) {
	        return new AppItem(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.id = source["id"];
	    }
	}
	export class AppCategory {
	    name: string;
	    apps: AppItem[];

	    static createFrom(source: any = {}) {
	        return new AppCategory(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.apps = this.convertValues(source["apps"], AppItem);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace sysinfo {

	export class BatteryStats {
	    status: string;
	    level: number;

	    static createFrom(source: any = {}) {
	        return new BatteryStats(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.level = source["level"];
	    }
	}
	export class NetworkStats {
	    status: string;
	    ping: string;

	    static createFrom(source: any = {}) {
	        return new NetworkStats(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.ping = source["ping"];
	    }
	}
	export class Partition {
	    name: string;
	    label: string;
	    total: number;
	    free: number;

	    static createFrom(source: any = {}) {
	        return new Partition(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.label = source["label"];
	        this.total = source["total"];
	        this.free = source["free"];
	    }
	}
	export class Info {
	    cpu: string;
	    cpu_usage: number;
	    memory: string;
	    mem_usage: number;
	    gpu: string;
	    disk: string;
	    partitions: Partition[];
	    os: string;
	    uptime: string;
	    disk_health: string;
	    network: NetworkStats;
	    battery: BatteryStats;

	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = source["cpu"];
	        this.cpu_usage = source["cpu_usage"];
	        this.memory = source["memory"];
	        this.mem_usage = source["mem_usage"];
	        this.gpu = source["gpu"];
	        this.disk = source["disk"];
	        this.partitions = this.convertValues(source["partitions"], Partition);
	        this.os = source["os"];
	        this.uptime = source["uptime"];
	        this.disk_health = source["disk_health"];
	        this.network = this.convertValues(source["network"], NetworkStats);
	        this.battery = this.convertValues(source["battery"], BatteryStats);
	    }

		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}


}
