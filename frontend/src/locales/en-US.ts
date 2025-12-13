// English locale messages
export default {
  // Common
  common: {
    add: "Add",
    edit: "Edit",
    delete: "Delete",
    save: "Save",
    cancel: "Cancel",
    confirm: "Confirm",
    refresh: "Refresh",
    loading: "Loading...",
    noData: "No data",
    today: "Today",
    yesterday: "Yesterday",
    total: "Total",
    status: "Status",
    actions: "Actions",
    configure: "Configure",
    next: "Next",
    prev: "Previous",
  },

  // Splash Screen
  splash: {
    tagline: "Streamline your freelance business",
    initializing: "Initializing...",
    start: "Start",
    welcomeBack: "Welcome back, entering your workspace...",
    progress: {
      backendInit: "Initializing backend",
      updateCheck: "Checking updates",
      sessionRestore: "Restoring session",
      preparingWorkspace: "Preparing workspace",
      ready: "Ready",
    },
  },

  // Authentication
  auth: {
    welcome: "Welcome Back",
    selectUser: "Select your profile to continue",
    addUser: "Add User",
    enterPassword: "Enter password",
    login: "Login",
    invalidPassword: "Invalid password",
    createAccount: "Create Account",
    setupProfile: "Set up your profile to get started",
    username: "Username",
    usernamePlaceholder: "Enter your username",
    email: "Email (Optional)",
    emailPlaceholder: "your@email.com",
    password: "Password",
    passwordPlaceholder: "Enter password",
    confirmPassword: "Confirm Password",
    confirmPasswordPlaceholder: "Re-enter password",
    passwordsNotMatch: "Passwords do not match",
    regenerateAvatar: "Generate new avatar",
    createProfile: "Create Profile",
    registerFailed: "Registration failed",
    logout: "Logout",
    switchUser: "Switch User",
    // Financial preferences
    financialPreferences: "Financial Preferences",
    language: "Language",
    currency: "Currency",
    province: "Province/Region (for tax)",
    timezone: "Timezone",
    // Step form
    stepProfile: "Profile",
    stepSecurity: "Security",
    stepPreferences: "Preferences",
    setPassword: "Set Password",
  },

  // Navigation
  nav: {
    dashboard: "Dashboard",
    clients: "Clients",
    projects: "Projects",
    timesheet: "Timesheet",
    invoices: "Invoices",
    reports: "Reports",
    finance: "Finance",
    settings: "Settings",
    help: "Help",
  },

  // Theme
  theme: {
    switchToLight: "Switch to Light Mode",
    switchToDark: "Switch to Dark Mode",
  },

  // Footer/Status Bar
  footer: {
    statusBar: "Status Bar:",
    weeklyHours: "Weekly Hours:",
    pendingPayment: "Pending Payment:",
  },

  // Timesheet
  timesheet: {
    title: "Timesheet",
    subtitle: "Track and log your billable hours",
    logTime: "Log Time",
    noEntries: "No time entries yet",
    noEntriesHint: "Start the timer above or add time manually",
    thisWeek: "This Week",
    // Timer
    timer: {
      placeholder: "What are you working on?",
      selectProject: "Select project",
      start: "Start",
      stop: "Stop",
      discard: "Discard timer",
      discardedMsg: "Timer discarded",
      loggedMsg: "Time logged successfully!",
      selectProjectFirst: "Please select a project first",
    },
    // Quick Entry
    quickEntry: {
      placeholder: "What did you work on?",
      project: "Project",
      date: "Date",
      duration: "Duration",
      durationHint: "Tip: Enter duration as {examples}",
      durationExamples: "1h, 1h 30m, 1.5, or 90m",
      invalidDuration:
        "Please enter a valid duration (e.g., 1h, 1h 30m, 1.5, 90m)",
      selectProject: "Please select a project",
      enterDescription: "Please enter a description",
      loggedMsg: "Time logged!",
    },
    // Entry
    entry: {
      noProject: "No Project",
      invoiced: "Invoiced",
      continueTimer: "Continue timer",
      editEntry: "Edit entry",
      deleteEntry: "Delete entry",
      deleteConfirm: "Delete this time entry?",
      deletedMsg: "Entry deleted",
      updatedMsg: "Time entry updated",
    },
    // Entries section
    entries: {
      title: "Time Entries",
      showingResults: "Showing {from} to {to} of {total} results",
      pageInfo: "Page {current} of {total}",
      exportCSV: "Export CSV",
      selectProject: "Select Project",
      describeTask: "Describe your task...",
      nonBillable: "Non-billable",
      billable: "Billable",
      addEntry: "Add Entry",
    },
    // Table columns
    columns: {
      date: "DATE",
      project: "PROJECT",
      task: "TASK",
      status: "STATUS",
      hours: "HOURS",
      billable: "BILLABLE",
    },
    // Form Modal
    form: {
      editTitle: "Edit Time Entry",
      createTitle: "Log Time",
      project: "Project",
      date: "Date",
      duration: "Duration (Hours)",
      description: "Description",
      descriptionPlaceholder: "What did you work on?",
      timeRange: "Time Range (Optional)",
      startTime: "Start time",
      endTime: "End time",
      alreadyInvoiced: "Already invoiced",
      update: "Update",
    },
  },

  // Clients
  clients: {
    title: "Clients",
    subtitle: "Manage your client relationships",
    addClient: "Add Client",
    editClient: "Edit Client",
    newClient: "New Client",
    noClients: "No clients found",
    deleteTitle: "Delete Client",
    deleteConfirm:
      'Are you sure you want to delete "{name}"? This action cannot be undone.',
    deleteSuccess: "Client deleted successfully",
    deleteError: "Failed to delete client",
    updateSuccess: "Client updated successfully",
    createSuccess: "Client created successfully",
    saveError: "Failed to save client",
    columns: {
      clientName: "Client Name",
      contactPerson: "Contact Person",
      status: "Status",
      actions: "Actions",
      associatedProjects: "Associated Projects",
      projectsCount: "{count} Projects",
    },
    status: {
      active: "Active",
      inactive: "Inactive",
    },
  },

  // Projects
  projects: {
    title: "Projects",
    subtitle: "Manage your projects",
    addProject: "Add Project",
    newProject: "New Project",
    editProject: "Edit Project",
    noProjects: "No projects found",
    deleteTitle: "Delete Project",
    deleteConfirm: 'Are you sure you want to delete "{name}"?',
    deleteSuccess: "Project deleted successfully",
    deleteError: "Failed to delete project",
    updateSuccess: "Project updated successfully",
    createSuccess: "Project created successfully",
    saveError: "Failed to save project",
    manage: "Managing {name}",
    columns: {
      projectName: "Project Name",
      status: "Status",
      progress: "Progress",
      deadline: "Deadline",
      noDeadline: "No Deadline",
      hourlyRate: "Hourly Rate",
      actions: "Actions",
    },
    status: {
      active: "Active",
      archived: "Archived",
      completed: "Completed",
    },
  },

  // Invoices
  invoices: {
    title: "Invoices",
    subtitle: "Manage your invoices",
    createInvoice: "Create Invoice",
    newInvoice: "New Invoice",
    editInvoice: "Edit Invoice",
    noInvoices: "No invoices found",
    updateSuccess: "Invoice updated successfully",
    createSuccess: "Invoice created successfully",
    saveError: "Failed to save invoice",
    downloading: "Downloading {number}...",
    stats: {
      outstandingAmount: "Outstanding Amount",
    },
    columns: {
      invoiceNumber: "Invoice #",
      client: "Client",
      issueDate: "Issue Date",
      amount: "Amount",
      status: "Status",
      actions: "Actions",
    },
    status: {
      paid: "Paid",
      sent: "Sent",
      overdue: "Overdue",
    },
  },

  // Forms
  form: {
    required: "Required",
    save: "Save",
    cancel: "Cancel",
    create: "Create",
    update: "Update",
    validation: {
      required: "Please enter {field}",
      email: "Please enter a valid email",
      number: "Please enter a valid number",
      select: "Please select {field}",
    },
    client: {
      name: "Client Name",
      namePlaceholder: "Enter client name",
      email: "Email",
      emailPlaceholder: "client@example.com",
      contactPerson: "Contact Person",
      contactPersonPlaceholder: "John Doe",
      website: "Website",
      websitePlaceholder: "https://example.com",
      address: "Address",
      addressPlaceholder: "Full address",
      currency: "Currency",
      status: "Status",
      notes: "Notes",
      notesPlaceholder: "Additional notes",
    },
    project: {
      client: "Client",
      clientPlaceholder: "Select client",
      name: "Project Name",
      namePlaceholder: "Enter project name",
      description: "Description",
      descriptionPlaceholder: "Project description",
      hourlyRate: "Hourly Rate",
      currency: "Currency",
      status: "Status",
      deadline: "Deadline",
      tags: "Tags",
    },
  },

  // Dashboard
  dashboard: {
    title: "Dashboard",
    subtitle: "Overview of your freelance business",
  },

  // Reports
  reports: {
    title: "Reports",
    subtitle: "Analytics and insights",
  },

  // Finance
  finance: {
    title: "Finance",
    subtitle: "Manage your accounts and track cash flow",
    nav: {
      overview: "Overview",
      accounts: "Accounts",
      transactions: "Transactions",
      import: "Import",
      categories: "Categories",
      reports: "Reports",
    },
    overview: {
      title: "Financial Overview",
      subtitle: "Your financial health at a glance",
      quickActions: "Quick Actions",
      viewAllAccounts: "View All Accounts",
      manageYourAccounts: "Manage your bank accounts and credit cards",
      viewAllTransactions: "View All Transactions",
      browseRecentTransactions: "Browse and manage your recent transactions",
      importStatement: "Import Bank Statement",
      importBankStatement: "Import transactions from bank CSV or OFX files",
      viewReports: "View Reports",
      analyzeYourFinance: "Analyze your income, expenses and trends",
    },
    summary: {
      totalBalance: "Total Balance",
      totalIncome: "Total Income",
      totalExpense: "Total Expenses",
      netCashFlow: "Net Cash Flow",
    },
    actions: {
      addTransaction: "Add Transaction",
      addAccount: "Add Account",
    },
    accounts: {
      title: "Accounts",
      addAccount: "Add Account",
      table: {
        name: "Account Name",
        type: "Type",
        currency: "Currency",
        balance: "Balance",
      },
      types: {
        checking: "Checking",
        savings: "Savings",
        credit: "Credit Card",
        investment: "Investment",
      },
      errors: {
        loadFailed: "Failed to load accounts",
      },
    },
    transactions: {
      title: "Transactions",
      addTransaction: "Add Transaction",
    },
    import: {
      title: "Import Bank Statement",
    },
    categories: {
      title: "Categories",
    },
    reports: {
      title: "Financial Reports",
    },
    filters: {
      dateRange: "Select date range",
      allAccounts: "All Accounts",
      allTypes: "All Types",
      income: "Income",
      expense: "Expense",
      transfer: "Transfer",
      adjustment: "Adjustment",
    },
    table: {
      date: "Date",
      account: "Account",
      type: "Type",
      description: "Description",
      category: "Category",
      amount: "Amount",
    },
    empty: {
      noTransactions: "No transactions found",
      addFirstTransaction: "Add First Transaction",
    },
  },

  // Settings
  settings: {
    index: {
      description: "Manage your application preferences and configuration",
    },
    general: {
      title: "General",
      description: "Language, currency, date format and other preferences",
    },
    profile: {
      title: "Profile",
      description: "Your personal information and company details",
    },
    invoice: {
      title: "Invoice",
      description: "Invoice defaults, terms and header information",
    },
    email: {
      title: "Email",
      description: "Email provider settings and templates",
    },
    finance: {
      title: "Finance",
      description: "Finance accounts and transaction settings",
    },
  },
};
