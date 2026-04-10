export namespace dto {
	
	export class AccountOutput {
	    id: number;
	    name: string;
	    type: string;
	    currency: string;
	    balance: number;
	    bankName: string;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new AccountOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.currency = source["currency"];
	        this.balance = source["balance"];
	        this.bankName = source["bankName"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class CategoryOutput {
	    id: number;
	    name: string;
	    type: string;
	    color: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.color = source["color"];
	        this.icon = source["icon"];
	    }
	}
	export class ChangePasswordInput {
	    id: number;
	    oldPassword: string;
	    newPassword: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangePasswordInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.oldPassword = source["oldPassword"];
	        this.newPassword = source["newPassword"];
	    }
	}
	export class ClientOutput {
	    id: number;
	    name: string;
	    email: string;
	    website: string;
	    avatar: string;
	    contactPerson: string;
	    address: string;
	    currency: string;
	    status: string;
	    notes: string;
	    billingCompany: string;
	    billingAddress: string;
	    billingCity: string;
	    billingProvince: string;
	    billingPostalCode: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.website = source["website"];
	        this.avatar = source["avatar"];
	        this.contactPerson = source["contactPerson"];
	        this.address = source["address"];
	        this.currency = source["currency"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.billingCompany = source["billingCompany"];
	        this.billingAddress = source["billingAddress"];
	        this.billingCity = source["billingCity"];
	        this.billingProvince = source["billingProvince"];
	        this.billingPostalCode = source["billingPostalCode"];
	    }
	}
	export class CreateAccountInput {
	    name: string;
	    type: string;
	    currency: string;
	    balance: number;
	    bankName: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateAccountInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.currency = source["currency"];
	        this.balance = source["balance"];
	        this.bankName = source["bankName"];
	    }
	}
	export class CreateCategoryInput {
	    name: string;
	    type: string;
	    color: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCategoryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.color = source["color"];
	        this.icon = source["icon"];
	    }
	}
	export class CreateClientInput {
	    name: string;
	    email: string;
	    website: string;
	    avatar: string;
	    contactPerson: string;
	    address: string;
	    currency: string;
	    status: string;
	    notes: string;
	    billingCompany: string;
	    billingAddress: string;
	    billingCity: string;
	    billingProvince: string;
	    billingPostalCode: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateClientInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.email = source["email"];
	        this.website = source["website"];
	        this.avatar = source["avatar"];
	        this.contactPerson = source["contactPerson"];
	        this.address = source["address"];
	        this.currency = source["currency"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.billingCompany = source["billingCompany"];
	        this.billingAddress = source["billingAddress"];
	        this.billingCity = source["billingCity"];
	        this.billingProvince = source["billingProvince"];
	        this.billingPostalCode = source["billingPostalCode"];
	    }
	}
	export class InvoiceItemInput {
	    description: string;
	    quantity: number;
	    unitPrice: number;
	    amount: number;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceItemInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.quantity = source["quantity"];
	        this.unitPrice = source["unitPrice"];
	        this.amount = source["amount"];
	    }
	}
	export class CreateInvoiceInput {
	    clientId: number;
	    number: string;
	    issueDate: string;
	    dueDate: string;
	    subtotal: number;
	    taxRate: number;
	    taxAmount: number;
	    total: number;
	    status: string;
	    items: InvoiceItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new CreateInvoiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clientId = source["clientId"];
	        this.number = source["number"];
	        this.issueDate = source["issueDate"];
	        this.dueDate = source["dueDate"];
	        this.subtotal = source["subtotal"];
	        this.taxRate = source["taxRate"];
	        this.taxAmount = source["taxAmount"];
	        this.total = source["total"];
	        this.status = source["status"];
	        this.items = this.convertValues(source["items"], InvoiceItemInput);
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
	export class CreateProjectInput {
	    clientId: number;
	    name: string;
	    description: string;
	    hourlyRate: number;
	    currency: string;
	    status: string;
	    deadline: string;
	    tags: string[];
	    serviceType: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.clientId = source["clientId"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.hourlyRate = source["hourlyRate"];
	        this.currency = source["currency"];
	        this.status = source["status"];
	        this.deadline = source["deadline"];
	        this.tags = source["tags"];
	        this.serviceType = source["serviceType"];
	    }
	}
	export class CreateTimeEntryInput {
	    projectId: number;
	    date: string;
	    startTime: string;
	    endTime: string;
	    durationSeconds: number;
	    description: string;
	    billable: boolean;
	    invoiced: boolean;
	    billingMode: string;
	    manualAmount?: number | null;
	
	    static createFrom(source: any = {}) {
	        return new CreateTimeEntryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectId = source["projectId"];
	        this.date = source["date"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.durationSeconds = source["durationSeconds"];
	        this.description = source["description"];
	        this.billable = source["billable"];
	        this.invoiced = source["invoiced"];
	        this.billingMode = source["billingMode"];
	        this.manualAmount = source["manualAmount"];
	    }
	}
	export class FinanceSummary {
	    totalBalance: number;
	    totalIncome: number;
	    totalExpense: number;
	    cashFlow: number;
	
	    static createFrom(source: any = {}) {
	        return new FinanceSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalBalance = source["totalBalance"];
	        this.totalIncome = source["totalIncome"];
	        this.totalExpense = source["totalExpense"];
	        this.cashFlow = source["cashFlow"];
	    }
	}
	export class ImportTransactionsInput {
	    accountId: number;
	    bankType: string;
	    fileContent: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportTransactionsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accountId = source["accountId"];
	        this.bankType = source["bankType"];
	        this.fileContent = source["fileContent"];
	    }
	}
	export class InvoiceEmailSettings {
	    provider: string;
	    from: string;
	    replyTo: string;
	    subjectTemplate: string;
	    bodyTemplate: string;
	    signature: string;
	    resendApiKey: string;
	    smtpHost: string;
	    smtpPort: number;
	    smtpUsername: string;
	    smtpPassword: string;
	    smtpUseTls: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceEmailSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.from = source["from"];
	        this.replyTo = source["replyTo"];
	        this.subjectTemplate = source["subjectTemplate"];
	        this.bodyTemplate = source["bodyTemplate"];
	        this.signature = source["signature"];
	        this.resendApiKey = source["resendApiKey"];
	        this.smtpHost = source["smtpHost"];
	        this.smtpPort = source["smtpPort"];
	        this.smtpUsername = source["smtpUsername"];
	        this.smtpPassword = source["smtpPassword"];
	        this.smtpUseTls = source["smtpUseTls"];
	    }
	}
	
	export class InvoiceItemOutput {
	    id: number;
	    description: string;
	    quantity: number;
	    unitPrice: number;
	    amount: number;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceItemOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.description = source["description"];
	        this.quantity = source["quantity"];
	        this.unitPrice = source["unitPrice"];
	        this.amount = source["amount"];
	    }
	}
	export class InvoiceOutput {
	    id: number;
	    clientId: number;
	    projectId: number;
	    number: string;
	    issueDate: string;
	    dueDate: string;
	    subtotal: number;
	    taxRate: number;
	    taxAmount: number;
	    total: number;
	    status: string;
	    items: InvoiceItemOutput[];
	
	    static createFrom(source: any = {}) {
	        return new InvoiceOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.clientId = source["clientId"];
	        this.projectId = source["projectId"];
	        this.number = source["number"];
	        this.issueDate = source["issueDate"];
	        this.dueDate = source["dueDate"];
	        this.subtotal = source["subtotal"];
	        this.taxRate = source["taxRate"];
	        this.taxAmount = source["taxAmount"];
	        this.total = source["total"];
	        this.status = source["status"];
	        this.items = this.convertValues(source["items"], InvoiceItemOutput);
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
	export class LoginInput {
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class ProjectOutput {
	    id: number;
	    clientId: number;
	    name: string;
	    description: string;
	    hourlyRate: number;
	    currency: string;
	    status: string;
	    deadline: string;
	    tags: string[];
	    serviceType: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.clientId = source["clientId"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.hourlyRate = source["hourlyRate"];
	        this.currency = source["currency"];
	        this.status = source["status"];
	        this.deadline = source["deadline"];
	        this.tags = source["tags"];
	        this.serviceType = source["serviceType"];
	    }
	}
	export class RegisterInput {
	    username: string;
	    password: string;
	    email: string;
	    avatarUrl: string;
	    settingsJson: string;
	
	    static createFrom(source: any = {}) {
	        return new RegisterInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	        this.email = source["email"];
	        this.avatarUrl = source["avatarUrl"];
	        this.settingsJson = source["settingsJson"];
	    }
	}
	export class ReportChartSeries {
	    dates: string[];
	    revenue: number[];
	    hours: number[];
	
	    static createFrom(source: any = {}) {
	        return new ReportChartSeries(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dates = source["dates"];
	        this.revenue = source["revenue"];
	        this.hours = source["hours"];
	    }
	}
	export class ReportFilter {
	    startDate: string;
	    endDate: string;
	    clientId: number;
	    projectId: number;
	
	    static createFrom(source: any = {}) {
	        return new ReportFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	        this.clientId = source["clientId"];
	        this.projectId = source["projectId"];
	    }
	}
	export class ReportRow {
	    date: string;
	    clientId: number;
	    clientName: string;
	    projectId: number;
	    projectName: string;
	    hours: number;
	    income: number;
	
	    static createFrom(source: any = {}) {
	        return new ReportRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.clientId = source["clientId"];
	        this.clientName = source["clientName"];
	        this.projectId = source["projectId"];
	        this.projectName = source["projectName"];
	        this.hours = source["hours"];
	        this.income = source["income"];
	    }
	}
	export class ReportOutput {
	    totalHours: number;
	    totalIncome: number;
	    rows: ReportRow[];
	    chart: ReportChartSeries;
	
	    static createFrom(source: any = {}) {
	        return new ReportOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalHours = source["totalHours"];
	        this.totalIncome = source["totalIncome"];
	        this.rows = this.convertValues(source["rows"], ReportRow);
	        this.chart = this.convertValues(source["chart"], ReportChartSeries);
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
	
	export class SetInvoiceTimeEntriesInput {
	    invoiceId: number;
	    timeEntryIds: number[];
	
	    static createFrom(source: any = {}) {
	        return new SetInvoiceTimeEntriesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.invoiceId = source["invoiceId"];
	        this.timeEntryIds = source["timeEntryIds"];
	    }
	}
	export class StatusBarOutput {
	    monthSeconds: number;
	    uninvoicedTotal: number;
	    unpaidTotal: number;
	    currency: string;
	
	    static createFrom(source: any = {}) {
	        return new StatusBarOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.monthSeconds = source["monthSeconds"];
	        this.uninvoicedTotal = source["uninvoicedTotal"];
	        this.unpaidTotal = source["unpaidTotal"];
	        this.currency = source["currency"];
	    }
	}
	export class TimeEntryOutput {
	    id: number;
	    projectId: number;
	    invoiceId: number;
	    date: string;
	    startTime: string;
	    endTime: string;
	    durationSeconds: number;
	    description: string;
	    billable: boolean;
	    invoiced: boolean;
	    billingMode: string;
	    manualAmount?: number;
	
	    static createFrom(source: any = {}) {
	        return new TimeEntryOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.invoiceId = source["invoiceId"];
	        this.date = source["date"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.durationSeconds = source["durationSeconds"];
	        this.description = source["description"];
	        this.billable = source["billable"];
	        this.invoiced = source["invoiced"];
	        this.billingMode = source["billingMode"];
	        this.manualAmount = source["manualAmount"];
	    }
	}
	export class TransactionFilter {
	    startDate?: string;
	    endDate?: string;
	    accountId?: number;
	
	    static createFrom(source: any = {}) {
	        return new TransactionFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
	        this.accountId = source["accountId"];
	    }
	}
	export class TransactionOutput {
	    id: number;
	    accountId: number;
	    categoryId?: number;
	    categoryName?: string;
	    categoryColor?: string;
	    // Go type: time
	    date: any;
	    description: string;
	    amount: number;
	    status: string;
	    referenceId: string;
	
	    static createFrom(source: any = {}) {
	        return new TransactionOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.accountId = source["accountId"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.categoryColor = source["categoryColor"];
	        this.date = this.convertValues(source["date"], null);
	        this.description = source["description"];
	        this.amount = source["amount"];
	        this.status = source["status"];
	        this.referenceId = source["referenceId"];
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
	export class UpdateAccountInput {
	    id: number;
	    name: string;
	    type: string;
	    currency: string;
	    balance: number;
	    bankName: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateAccountInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.currency = source["currency"];
	        this.balance = source["balance"];
	        this.bankName = source["bankName"];
	    }
	}
	export class UpdateCategoryInput {
	    id: number;
	    name: string;
	    type: string;
	    color: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCategoryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.color = source["color"];
	        this.icon = source["icon"];
	    }
	}
	export class UpdateClientInput {
	    id: number;
	    name: string;
	    email: string;
	    website: string;
	    avatar: string;
	    contactPerson: string;
	    address: string;
	    currency: string;
	    status: string;
	    notes: string;
	    billingCompany: string;
	    billingAddress: string;
	    billingCity: string;
	    billingProvince: string;
	    billingPostalCode: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateClientInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.website = source["website"];
	        this.avatar = source["avatar"];
	        this.contactPerson = source["contactPerson"];
	        this.address = source["address"];
	        this.currency = source["currency"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.billingCompany = source["billingCompany"];
	        this.billingAddress = source["billingAddress"];
	        this.billingCity = source["billingCity"];
	        this.billingProvince = source["billingProvince"];
	        this.billingPostalCode = source["billingPostalCode"];
	    }
	}
	export class UpdateInvoiceInput {
	    id: number;
	    clientId: number;
	    number: string;
	    issueDate: string;
	    dueDate: string;
	    subtotal: number;
	    taxRate: number;
	    taxAmount: number;
	    total: number;
	    status: string;
	    items: InvoiceItemInput[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateInvoiceInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.clientId = source["clientId"];
	        this.number = source["number"];
	        this.issueDate = source["issueDate"];
	        this.dueDate = source["dueDate"];
	        this.subtotal = source["subtotal"];
	        this.taxRate = source["taxRate"];
	        this.taxAmount = source["taxAmount"];
	        this.total = source["total"];
	        this.status = source["status"];
	        this.items = this.convertValues(source["items"], InvoiceItemInput);
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
	export class UpdateProjectInput {
	    id: number;
	    clientId: number;
	    name: string;
	    description: string;
	    hourlyRate: number;
	    currency: string;
	    status: string;
	    deadline: string;
	    tags: string[];
	    serviceType: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProjectInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.clientId = source["clientId"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.hourlyRate = source["hourlyRate"];
	        this.currency = source["currency"];
	        this.status = source["status"];
	        this.deadline = source["deadline"];
	        this.tags = source["tags"];
	        this.serviceType = source["serviceType"];
	    }
	}
	export class UpdateTimeEntryInput {
	    id: number;
	    projectId: number;
	    invoiceId: number;
	    date: string;
	    startTime: string;
	    endTime: string;
	    durationSeconds: number;
	    description: string;
	    billable: boolean;
	    invoiced: boolean;
	    billingMode: string;
	    manualAmount?: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateTimeEntryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.invoiceId = source["invoiceId"];
	        this.date = source["date"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.durationSeconds = source["durationSeconds"];
	        this.description = source["description"];
	        this.billable = source["billable"];
	        this.invoiced = source["invoiced"];
	        this.billingMode = source["billingMode"];
	        this.manualAmount = source["manualAmount"];
	    }
	}
	export class UpdateUserInput {
	    id: number;
	    username: string;
	    email: string;
	    avatarUrl: string;
	    settingsJson: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateUserInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.email = source["email"];
	        this.avatarUrl = source["avatarUrl"];
	        this.settingsJson = source["settingsJson"];
	    }
	}
	export class UserInvoiceSettings {
	    userId: number;
	    senderName: string;
	    senderCompany: string;
	    senderAddress: string;
	    senderPhone: string;
	    senderEmail: string;
	    senderPostalCode: string;
	    defaultTerms: string;
	    defaultMessageTemplate: string;
	
	    static createFrom(source: any = {}) {
	        return new UserInvoiceSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.senderName = source["senderName"];
	        this.senderCompany = source["senderCompany"];
	        this.senderAddress = source["senderAddress"];
	        this.senderPhone = source["senderPhone"];
	        this.senderEmail = source["senderEmail"];
	        this.senderPostalCode = source["senderPostalCode"];
	        this.defaultTerms = source["defaultTerms"];
	        this.defaultMessageTemplate = source["defaultMessageTemplate"];
	    }
	}
	export class UserListItem {
	    id: number;
	    username: string;
	    avatarUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new UserListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.avatarUrl = source["avatarUrl"];
	    }
	}
	export class UserOutput {
	    id: number;
	    uuid: string;
	    username: string;
	    email: string;
	    avatarUrl: string;
	    createdAt: string;
	    lastLogin: string;
	    settingsJson: string;
	
	    static createFrom(source: any = {}) {
	        return new UserOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.uuid = source["uuid"];
	        this.username = source["username"];
	        this.email = source["email"];
	        this.avatarUrl = source["avatarUrl"];
	        this.createdAt = source["createdAt"];
	        this.lastLogin = source["lastLogin"];
	        this.settingsJson = source["settingsJson"];
	    }
	}
	export class UserPreferences {
	    userId: number;
	    currency: string;
	    language: string;
	    theme: string;
	    timezone: string;
	    dateFormat: string;
	    moduleOverrides?: Record<string, boolean>;
	
	    static createFrom(source: any = {}) {
	        return new UserPreferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.currency = source["currency"];
	        this.language = source["language"];
	        this.theme = source["theme"];
	        this.timezone = source["timezone"];
	        this.dateFormat = source["dateFormat"];
	        this.moduleOverrides = source["moduleOverrides"];
	    }
	}
	export class UserSettings {
	    currency: string;
	    defaultTaxRate: number;
	    language: string;
	    theme: string;
	    dateFormat: string;
	    timezone: string;
	    senderName: string;
	    senderCompany: string;
	    senderAddress: string;
	    senderPhone: string;
	    senderEmail: string;
	    senderPostalCode: string;
	    invoiceTerms: string;
	    defaultMessageTemplate: string;
	    moduleOverrides?: Record<string, boolean>;
	    hstRegistered: boolean;
	    hstNumber: string;
	    taxEnabled: boolean;
	    expectedIncome: string;
	
	    static createFrom(source: any = {}) {
	        return new UserSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currency = source["currency"];
	        this.defaultTaxRate = source["defaultTaxRate"];
	        this.language = source["language"];
	        this.theme = source["theme"];
	        this.dateFormat = source["dateFormat"];
	        this.timezone = source["timezone"];
	        this.senderName = source["senderName"];
	        this.senderCompany = source["senderCompany"];
	        this.senderAddress = source["senderAddress"];
	        this.senderPhone = source["senderPhone"];
	        this.senderEmail = source["senderEmail"];
	        this.senderPostalCode = source["senderPostalCode"];
	        this.invoiceTerms = source["invoiceTerms"];
	        this.defaultMessageTemplate = source["defaultMessageTemplate"];
	        this.moduleOverrides = source["moduleOverrides"];
	        this.hstRegistered = source["hstRegistered"];
	        this.hstNumber = source["hstNumber"];
	        this.taxEnabled = source["taxEnabled"];
	        this.expectedIncome = source["expectedIncome"];
	    }
	}
	export class UserTaxSettings {
	    userId: number;
	    hstRegistered: boolean;
	    hstNumber: string;
	    taxEnabled: boolean;
	    defaultTaxRate: number;
	    expectedIncome: string;
	
	    static createFrom(source: any = {}) {
	        return new UserTaxSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.hstRegistered = source["hstRegistered"];
	        this.hstNumber = source["hstNumber"];
	        this.taxEnabled = source["taxEnabled"];
	        this.defaultTaxRate = source["defaultTaxRate"];
	        this.expectedIncome = source["expectedIncome"];
	    }
	}

}

export namespace main {
	
	export class BootTimings {
	    // Go type: time
	    processStart: any;
	    dbInitMs: number;
	    servicesInitMs: number;
	    totalBeforeUiMs: number;
	
	    static createFrom(source: any = {}) {
	        return new BootTimings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.processStart = this.convertValues(source["processStart"], null);
	        this.dbInitMs = source["dbInitMs"];
	        this.servicesInitMs = source["servicesInitMs"];
	        this.totalBeforeUiMs = source["totalBeforeUiMs"];
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

export namespace update {
	
	export class Platform {
	    url: string;
	    signature: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new Platform(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.signature = source["signature"];
	        this.size = source["size"];
	    }
	}
	export class Info {
	    version: string;
	    // Go type: time
	    releaseDate: any;
	    releaseNotes: string;
	    releaseNotesUrl?: string;
	    mandatory: boolean;
	    minimumOsVersion?: Record<string, string>;
	    platforms: Record<string, Platform>;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.releaseDate = this.convertValues(source["releaseDate"], null);
	        this.releaseNotes = source["releaseNotes"];
	        this.releaseNotesUrl = source["releaseNotesUrl"];
	        this.mandatory = source["mandatory"];
	        this.minimumOsVersion = source["minimumOsVersion"];
	        this.platforms = this.convertValues(source["platforms"], Platform, true);
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
	
	export class State {
	    status: string;
	    currentVersion: string;
	    latestVersion?: string;
	    updateInfo?: Info;
	    downloadProgress?: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new State(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.updateInfo = this.convertValues(source["updateInfo"], Info);
	        this.downloadProgress = source["downloadProgress"];
	        this.error = source["error"];
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

