export namespace catalog {
	
	export class Repo {
	    id: number;
	    folderId: number;
	    folderPath: string;
	    name: string;
	    path: string;
	    category: string;
	    readmePath: string;
	    readmeText: string;
	    fileCount: number;
	    scriptCount: number;
	
	    static createFrom(source: any = {}) {
	        return new Repo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.folderId = source["folderId"];
	        this.folderPath = source["folderPath"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.category = source["category"];
	        this.readmePath = source["readmePath"];
	        this.readmeText = source["readmeText"];
	        this.fileCount = source["fileCount"];
	        this.scriptCount = source["scriptCount"];
	    }
	}
	export class Category {
	    name: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new Category(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.count = source["count"];
	    }
	}
	export class Folder {
	    id: number;
	    path: string;
	    displayName: string;
	    isDefault: boolean;
	    lastScanned: string;
	    repoCount: number;
	    scriptCount: number;
	
	    static createFrom(source: any = {}) {
	        return new Folder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.displayName = source["displayName"];
	        this.isDefault = source["isDefault"];
	        this.lastScanned = source["lastScanned"];
	        this.repoCount = source["repoCount"];
	        this.scriptCount = source["scriptCount"];
	    }
	}
	export class CatalogView {
	    folders: Folder[];
	    categories: Category[];
	    repos: Repo[];
	    totalScripts: number;
	
	    static createFrom(source: any = {}) {
	        return new CatalogView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.folders = this.convertValues(source["folders"], Folder);
	        this.categories = this.convertValues(source["categories"], Category);
	        this.repos = this.convertValues(source["repos"], Repo);
	        this.totalScripts = source["totalScripts"];
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
	
	export class CustomCategory {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CustomCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class CustomRepo {
	    id: number;
	    name: string;
	    url: string;
	    category: string;
	    summary: string;
	
	    static createFrom(source: any = {}) {
	        return new CustomRepo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.url = source["url"];
	        this.category = source["category"];
	        this.summary = source["summary"];
	    }
	}
	export class CustomSoftware {
	    id: number;
	    name: string;
	    version: string;
	    category: string;
	    download: string;
	    notes: string;
	    wingetId: string;
	
	    static createFrom(source: any = {}) {
	        return new CustomSoftware(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.version = source["version"];
	        this.category = source["category"];
	        this.download = source["download"];
	        this.notes = source["notes"];
	        this.wingetId = source["wingetId"];
	    }
	}
	export class ExportResult {
	    copied: number;
	    skipped: number;
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new ExportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.copied = source["copied"];
	        this.skipped = source["skipped"];
	        this.errors = source["errors"];
	    }
	}
	
	
	export class ScriptFile {
	    id: number;
	    repoId: number;
	    repo: string;
	    relPath: string;
	    absPath: string;
	    name: string;
	    ext: string;
	    lang: string;
	    size: number;
	    mtime: number;
	    isScript: boolean;
	    snippet: string;
	
	    static createFrom(source: any = {}) {
	        return new ScriptFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.repoId = source["repoId"];
	        this.repo = source["repo"];
	        this.relPath = source["relPath"];
	        this.absPath = source["absPath"];
	        this.name = source["name"];
	        this.ext = source["ext"];
	        this.lang = source["lang"];
	        this.size = source["size"];
	        this.mtime = source["mtime"];
	        this.isScript = source["isScript"];
	        this.snippet = source["snippet"];
	    }
	}
	export class RepoDetail {
	    repo?: Repo;
	    scripts: ScriptFile[];
	
	    static createFrom(source: any = {}) {
	        return new RepoDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.repo = this.convertValues(source["repo"], Repo);
	        this.scripts = this.convertValues(source["scripts"], ScriptFile);
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
	export class RunResult {
	    command: string;
	    launched: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new RunResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.command = source["command"];
	        this.launched = source["launched"];
	        this.message = source["message"];
	    }
	}
	
	export class ScriptRow {
	    id: number;
	    name: string;
	    category: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ScriptRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.size = source["size"];
	    }
	}
	export class SearchResult {
	    repos: Repo[];
	    scripts: ScriptFile[];
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.repos = this.convertValues(source["repos"], Repo);
	        this.scripts = this.convertValues(source["scripts"], ScriptFile);
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

export namespace githublib {
	
	export class InstallResult {
	    name: string;
	    status: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.status = source["status"];
	        this.message = source["message"];
	    }
	}
	export class RecommendedRepo {
	    name: string;
	    category: string;
	    summary: string;
	
	    static createFrom(source: any = {}) {
	        return new RecommendedRepo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.category = source["category"];
	        this.summary = source["summary"];
	    }
	}
	export class RepoInfo {
	    name: string;
	    localPath: string;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new RepoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.localPath = source["localPath"];
	        this.category = source["category"];
	    }
	}
	export class SoftwareItem {
	    name: string;
	    version: string;
	    category: string;
	    download: string;
	    notes: string;
	    wingetId: string;
	
	    static createFrom(source: any = {}) {
	        return new SoftwareItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.category = source["category"];
	        this.download = source["download"];
	        this.notes = source["notes"];
	        this.wingetId = source["wingetId"];
	    }
	}

}

export namespace main {
	
	export class LocalItem {
	    name: string;
	    path: string;
	    isDir: boolean;
	    size: number;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new LocalItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	        this.kind = source["kind"];
	    }
	}

}

export namespace update {
	
	export class Info {
	    available: boolean;
	    version: string;
	    build: number;
	    installer_url: string;
	    notes: string;
	    force_update: boolean;
	    sha256?: string;
	    stale: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.version = source["version"];
	        this.build = source["build"];
	        this.installer_url = source["installer_url"];
	        this.notes = source["notes"];
	        this.force_update = source["force_update"];
	        this.sha256 = source["sha256"];
	        this.stale = source["stale"];
	    }
	}

}

