import type {
  IClientService,
  IProjectService,
  ITimeEntryService,
  IInvoiceService,
  Client,
  Project,
  TimeEntry,
  Invoice,
  CreateClientInput,
  UpdateClientInput,
  CreateProjectInput,
  UpdateProjectInput,
  CreateTimeEntryInput,
  UpdateTimeEntryInput,
  CreateInvoiceInput,
  UpdateInvoiceInput,
} from "@/types";
import { dto } from "@/wailsjs/go/models";

// Mock Data
let clients: Client[] = [
  {
    id: 1,
    name: "TechCorp Inc.",
    email: "billing@techcorp.com",
    website: "https://techcorp.com",
    contactPerson: "Sarah Jenkins",
    address: "123 Tech Blvd, Toronto",
    currency: "CAD",
    status: "active",
    avatar: "",
    notes: "",
  },
  {
    id: 2,
    name: "StartupStudio",
    email: "hello@startup.io",
    website: "https://startup.io",
    contactPerson: "Mike Ross",
    address: "456 Innovation Dr, Waterloo",
    currency: "USD",
    status: "active",
    avatar: "",
    notes: "",
  },
  {
    id: 3,
    name: "Legacy Systems",
    email: "accounts@legacy.net",
    website: "",
    currency: "CAD",
    status: "inactive",
    contactPerson: "Old Greg",
    address: "",
    avatar: "",
    notes: "",
  },
];

let projects: Project[] = [
  {
    id: 1,
    clientId: 1,
    name: "SaaS Platform Revamp",
    description: "Modernizing the legacy CRM with Vue 3 and Node.js.",
    hourlyRate: 120,
    currency: "CAD",
    status: "active",
    deadline: "2025-10-31",
    tags: ["Vue", "Node", "Architecture"],
  },
  {
    id: 2,
    clientId: 1,
    name: "Mobile App API",
    description: "RESTful API development for the new iOS app.",
    hourlyRate: 120,
    currency: "CAD",
    status: "completed",
    deadline: "2025-04-15",
    tags: ["API", "Express"],
  },
  {
    id: 3,
    clientId: 2,
    name: "MVP Development",
    description: "Rapid prototyping of the initial product concept.",
    hourlyRate: 95,
    currency: "USD",
    status: "active",
    deadline: "2025-08-01",
    tags: ["Prototype", "FullStack"],
  },
];

let timeEntries: TimeEntry[] = [
  {
    id: 1,
    projectId: 1,
    date: "2025-05-12",
    startTime: "09:00",
    endTime: "13:00",
    durationSeconds: 14400,
    description: "Database schema design and migration planning",
    invoiced: true,
  },
  {
    id: 2,
    projectId: 1,
    date: "2025-05-12",
    startTime: "14:00",
    endTime: "18:00",
    durationSeconds: 14400,
    description: "API implementation - Auth Service",
    invoiced: true,
  },
  {
    id: 3,
    projectId: 3,
    date: "2025-05-13",
    startTime: "10:00",
    endTime: "13:00",
    durationSeconds: 10800,
    description: "Frontend setup - Vite + Naive UI",
    invoiced: false,
  },
  {
    id: 4,
    projectId: 2,
    date: "2025-05-13",
    startTime: "14:00",
    endTime: "16:30",
    durationSeconds: 9000,
    description: "Mobile API Integration testing",
    invoiced: false,
  },
  {
    id: 5,
    projectId: 1,
    date: "2025-05-14",
    startTime: "",
    endTime: "",
    durationSeconds: 7200,
    description: "Client meeting & Requirements gathering",
    invoiced: false,
  },
];

let invoices: Invoice[] = [
  dto.InvoiceOutput.createFrom({
    id: 1,
    clientId: 1,
    number: "INV-2025-001",
    issueDate: "2025-05-01",
    dueDate: "2025-05-15",
    status: "paid",
    subtotal: 1200,
    taxRate: 0.13,
    taxAmount: 156,
    total: 1356,
    items: [
      {
        id: 101,
        description: "Consulting Services - April",
        quantity: 10,
        unitPrice: 120,
        amount: 1200,
      },
    ],
  }),
  dto.InvoiceOutput.createFrom({
    id: 2,
    clientId: 2,
    number: "INV-2025-002",
    issueDate: "2025-05-10",
    dueDate: "2025-05-24",
    status: "sent",
    subtotal: 2850,
    taxRate: 0.0,
    taxAmount: 0,
    total: 2850,
    items: [
      {
        id: 102,
        description: "MVP Phase 1 Development",
        quantity: 30,
        unitPrice: 95,
        amount: 2850,
      },
    ],
  }),
  dto.InvoiceOutput.createFrom({
    id: 3,
    clientId: 1,
    number: "INV-2025-003",
    issueDate: "2025-05-14",
    dueDate: "2025-05-28",
    status: "draft",
    subtotal: 0,
    taxRate: 0.13,
    taxAmount: 0,
    total: 0,
    items: [],
  }),
];

// Mock Services
export const mockClientService: IClientService = {
  async list() {
    return [...clients];
  },
  async get(id) {
    const client = clients.find((c) => c.id === id);
    if (!client) throw new Error("Client not found");
    return client;
  },
  async create(input: CreateClientInput) {
    const newClient: Client = {
      ...input,
      id: Math.max(0, ...clients.map((c) => c.id)) + 1,
    };
    clients.push(newClient);
    return newClient;
  },
  async update(input: UpdateClientInput): Promise<Client> {
    const index = clients.findIndex((c) => c.id === input.id);
    if (index === -1) throw new Error("Client not found");
    clients[index] = { ...clients[index], ...input };
    return clients[index];
  },
  async delete(id) {
    clients = clients.filter((c) => c.id !== id);
  },
};

export const mockProjectService: IProjectService = {
  async list() {
    return [...projects];
  },
  async listByClient(clientId) {
    return projects.filter((p) => p.clientId === clientId);
  },
  async get(id) {
    const project = projects.find((p) => p.id === id);
    if (!project) throw new Error("Project not found");
    return project;
  },
  async create(input: CreateProjectInput) {
    const newProject: Project = {
      ...input,
      id: Math.max(0, ...projects.map((p) => p.id)) + 1,
      description: input.description ?? "",
      deadline: input.deadline ?? "",
      tags: input.tags ?? [],
    };
    projects.push(newProject);
    return newProject;
  },
  async update(input: UpdateProjectInput): Promise<Project> {
    const index = projects.findIndex((p) => p.id === input.id);
    if (index === -1) throw new Error("Project not found");
    projects[index] = {
      ...projects[index],
      ...input,
      description: input.description ?? "",
      deadline: input.deadline ?? "",
      tags: input.tags ?? [],
    };
    return projects[index];
  },
  async delete(id) {
    projects = projects.filter((p) => p.id !== id);
  },
};

export const mockTimeEntryService: ITimeEntryService = {
  async list(projectId?: number) {
    return timeEntries.filter((t) => {
      if (projectId && t.projectId !== projectId) return false;
      return true;
    });
  },
  async get(id) {
    const entry = timeEntries.find((t) => t.id === id);
    if (!entry) throw new Error("Time entry not found");
    return entry;
  },
  async create(input: CreateTimeEntryInput) {
    const newEntry: TimeEntry = {
      ...input,
      id: Math.max(0, ...timeEntries.map((t) => t.id)) + 1,
      startTime: input.startTime ?? "",
      endTime: input.endTime ?? "",
      invoiced: input.invoiced ?? false,
    };
    timeEntries.push(newEntry);
    return newEntry;
  },
  async update(input: UpdateTimeEntryInput): Promise<TimeEntry> {
    const index = timeEntries.findIndex((t) => t.id === input.id);
    if (index === -1) throw new Error("Time entry not found");
    timeEntries[index] = {
      ...timeEntries[index],
      ...input,
      startTime: input.startTime ?? "",
      endTime: input.endTime ?? "",
      invoiced: input.invoiced ?? false,
    };
    return timeEntries[index];
  },
  async delete(id) {
    timeEntries = timeEntries.filter((t) => t.id !== id);
  },
};

export const mockInvoiceService: IInvoiceService = {
  async list() {
    return [...invoices];
  },
  async get(id) {
    const invoice = invoices.find((i) => i.id === id);
    if (!invoice) throw new Error("Invoice not found");
    return invoice;
  },
  async create(input: CreateInvoiceInput) {
    const newInvoice: Invoice = dto.InvoiceOutput.createFrom({
      ...input,
      id: Math.max(0, ...invoices.map((i) => i.id)) + 1,
      items: input.items.map((item, idx) => ({
        ...item,
        id:
          Math.max(0, ...invoices.flatMap((i) => i.items.map((it) => it.id))) +
          idx +
          1,
      })),
    });
    invoices.push(newInvoice);
    return newInvoice;
  },
  async update(input: UpdateInvoiceInput): Promise<Invoice> {
    const index = invoices.findIndex((i) => i.id === input.id);
    if (index === -1) throw new Error("Invoice not found");
    const invoice = invoices[index]!;
    const existingItems = invoice.items;
    invoices[index] = dto.InvoiceOutput.createFrom({
      ...invoice,
      ...input,
      items: input.items.map((item, idx) => ({
        ...item,
        id: existingItems[idx]?.id ?? idx + 1,
      })),
    });
    return invoices[index]!;
  },
  async delete(id) {
    invoices = invoices.filter((i) => i.id !== id);
  },
  async generatePdf(_id?: number, _message?: string) {
    return "mock-pdf-base64";
  },
  async sendEmail() {
    return true;
  },
  async setTimeEntries(input: dto.SetInvoiceTimeEntriesInput) {
    // Mock: no-op, return first invoice
    return invoices.find((i) => i.id === input.invoiceId) ?? invoices[0]!;
  },
};
