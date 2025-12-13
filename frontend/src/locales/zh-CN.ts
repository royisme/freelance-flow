// 简体中文 locale messages
export default {
  // 通用
  common: {
    add: "添加",
    edit: "编辑",
    delete: "删除",
    save: "保存",
    cancel: "取消",
    confirm: "确认",
    refresh: "刷新",
    loading: "加载中...",
    noData: "暂无数据",
    today: "今天",
    yesterday: "昨天",
    total: "合计",
    status: "状态",
    actions: "操作",
    configure: "配置",
    next: "下一步",
    prev: "上一步",
  },

  // 启动画面
  splash: {
    tagline: "简化您的自由职业业务管理",
    initializing: "正在初始化...",
    start: "启动",
    welcomeBack: "欢迎回来，正在进入工作台...",
    progress: {
      backendInit: "正在初始化后端",
      updateCheck: "正在检查更新",
      sessionRestore: "正在恢复会话",
      preparingWorkspace: "正在准备工作台",
      ready: "准备完成",
    },
  },

  // 认证
  auth: {
    welcome: "欢迎回来",
    selectUser: "选择您的账户继续",
    addUser: "添加用户",
    enterPassword: "输入密码",
    login: "登录",
    invalidPassword: "密码错误",
    createAccount: "创建账户",
    setupProfile: "设置您的个人资料开始使用",
    username: "用户名",
    usernamePlaceholder: "请输入用户名",
    email: "邮箱（可选）",
    emailPlaceholder: "your@email.com",
    password: "密码",
    passwordPlaceholder: "请输入密码",
    confirmPassword: "确认密码",
    confirmPasswordPlaceholder: "再次输入密码",
    passwordsNotMatch: "两次输入的密码不一致",
    regenerateAvatar: "重新生成头像",
    createProfile: "创建账户",
    registerFailed: "注册失败",
    logout: "退出登录",
    switchUser: "切换用户",
    // 财务偏好
    financialPreferences: "财务偏好设置",
    language: "界面语言",
    currency: "默认货币",
    province: "省份/地区（用于税率）",
    timezone: "时区",
    // 分步表单
    stepProfile: "资料",
    stepSecurity: "安全",
    stepPreferences: "偏好",
    setPassword: "设置密码",
  },

  // 导航
  nav: {
    dashboard: "仪表盘",
    clients: "客户",
    projects: "项目",
    timesheet: "工时",
    invoices: "发票",
    reports: "报表",
    finance: "财务",
    settings: "设置",
    help: "帮助",
  },

  // 主题
  theme: {
    switchToLight: "切换到浅色模式",
    switchToDark: "切换到深色模式",
  },

  // 底部状态栏
  footer: {
    statusBar: "状态栏：",
    weeklyHours: "本周工时：",
    pendingPayment: "待收款：",
  },

  // 工时
  timesheet: {
    title: "工时",
    subtitle: "记录和跟踪您的计费工时",
    logTime: "记录工时",
    noEntries: "暂无工时记录",
    noEntriesHint: "使用上方计时器或手动添加工时",
    thisWeek: "本周",
    // 计时器
    timer: {
      placeholder: "正在进行什么工作？",
      selectProject: "选择项目",
      start: "开始",
      stop: "停止",
      discard: "丢弃计时",
      discardedMsg: "计时已丢弃",
      loggedMsg: "工时记录成功！",
      selectProjectFirst: "请先选择一个项目",
    },
    // 快速录入
    quickEntry: {
      placeholder: "您完成了什么工作？",
      project: "项目",
      date: "日期",
      duration: "时长",
      durationHint: "提示：时长格式支持 {examples}",
      durationExamples: "1h、1h 30m、1.5 或 90m",
      invalidDuration: "请输入有效的时长格式（如 1h、1h 30m、1.5、90m）",
      selectProject: "请选择项目",
      enterDescription: "请输入工作描述",
      loggedMsg: "工时已记录！",
    },
    // 条目
    entry: {
      noProject: "无项目",
      invoiced: "已开票",
      continueTimer: "继续计时",
      editEntry: "编辑条目",
      deleteEntry: "删除条目",
      deleteConfirm: "确定删除此工时记录？",
      deletedMsg: "条目已删除",
      updatedMsg: "工时记录已更新",
    },
    // 条目区域
    entries: {
      title: "工时记录",
      showingResults: "显示 {from} 到 {to} 条，共 {total} 条",
      pageInfo: "第 {current} 页，共 {total} 页",
      exportCSV: "导出 CSV",
      selectProject: "选择项目",
      describeTask: "描述您的任务...",
      nonBillable: "不可计费",
      billable: "可计费",
      addEntry: "添加记录",
    },
    // 表格列
    columns: {
      date: "日期",
      project: "项目",
      task: "任务",
      status: "状态",
      hours: "时长",
      billable: "计费金额",
    },
    // 表单弹窗
    form: {
      editTitle: "编辑工时记录",
      createTitle: "记录工时",
      project: "项目",
      date: "日期",
      duration: "时长（小时）",
      description: "工作描述",
      descriptionPlaceholder: "您完成了什么工作？",
      timeRange: "时间段（可选）",
      startTime: "开始时间",
      endTime: "结束时间",
      alreadyInvoiced: "已开票",
      update: "更新",
    },
  },

  // 客户
  clients: {
    title: "客户",
    subtitle: "管理您的客户关系",
    addClient: "添加客户",
    editClient: "编辑客户",
    newClient: "新建客户",
    noClients: "暂无客户",
    deleteTitle: "删除客户",
    deleteConfirm: '确定要删除 "{name}" 吗？此操作无法撤销。',
    deleteSuccess: "客户删除成功",
    deleteError: "删除客户失败",
    updateSuccess: "客户更新成功",
    createSuccess: "客户创建成功",
    saveError: "保存客户失败",
    columns: {
      clientName: "客户名称",
      contactPerson: "联系人",
      status: "状态",
      actions: "操作",
      associatedProjects: "关联项目",
      projectsCount: "{count} 个项目",
    },
    status: {
      active: "活跃",
      inactive: "非活跃",
    },
  },

  // 项目
  projects: {
    title: "项目",
    subtitle: "管理您的项目",
    addProject: "添加项目",
    newProject: "新建项目",
    editProject: "编辑项目",
    noProjects: "暂无项目",
    deleteTitle: "删除项目",
    deleteConfirm: '确定要删除 "{name}" 吗？',
    deleteSuccess: "项目删除成功",
    deleteError: "删除项目失败",
    updateSuccess: "项目更新成功",
    createSuccess: "项目创建成功",
    saveError: "保存项目失败",
    manage: "管理 {name}",
    columns: {
      projectName: "项目名称",
      status: "状态",
      progress: "进度",
      deadline: "截止日期",
      noDeadline: "无截止日期",
      hourlyRate: "时薪",
      actions: "操作",
    },
    status: {
      active: "进行中",
      archived: "已归档",
      completed: "已完成",
    },
  },

  // 发票
  invoices: {
    title: "发票",
    subtitle: "管理您的发票",
    createInvoice: "创建发票",
    newInvoice: "新建发票",
    editInvoice: "编辑发票",
    noInvoices: "暂无发票",
    updateSuccess: "发票更新成功",
    createSuccess: "发票创建成功",
    saveError: "保存发票失败",
    downloading: "正在下载 {number}...",
    stats: {
      outstandingAmount: "待收金额",
    },
    columns: {
      invoiceNumber: "发票号",
      client: "客户",
      issueDate: "开票日期",
      amount: "金额",
      status: "状态",
      actions: "操作",
    },
    status: {
      paid: "已支付",
      sent: "已发送",
      overdue: "已逾期",
    },
  },

  // 表单
  form: {
    required: "必填",
    save: "保存",
    cancel: "取消",
    create: "创建",
    update: "更新",
    validation: {
      required: "请输入{field}",
      email: "请输入有效的邮箱地址",
      number: "请输入有效的数字",
      select: "请选择{field}",
    },
    client: {
      name: "客户名称",
      namePlaceholder: "输入客户名称",
      email: "邮箱",
      emailPlaceholder: "client@example.com",
      contactPerson: "联系人",
      contactPersonPlaceholder: "张三",
      website: "网站",
      websitePlaceholder: "https://example.com",
      address: "地址",
      addressPlaceholder: "详细地址",
      currency: "货币",
      status: "状态",
      notes: "备注",
      notesPlaceholder: "额外备注信息",
    },
    project: {
      client: "客户",
      clientPlaceholder: "选择客户",
      name: "项目名称",
      namePlaceholder: "输入项目名称",
      description: "描述",
      descriptionPlaceholder: "项目描述",
      hourlyRate: "时薪",
      currency: "货币",
      status: "状态",
      deadline: "截止日期",
      tags: "标签",
    },
  },

  // 仪表盘
  dashboard: {
    title: "仪表盘",
    subtitle: "您的自由职业业务概览",
  },

  // 报表
  reports: {
    title: "报表",
    subtitle: "数据分析与洞察",
  },

  // 财务
  finance: {
    title: "财务",
    subtitle: "管理您的账户并跟踪现金流",
    nav: {
      overview: "概览",
      accounts: "账户",
      transactions: "交易",
      import: "导入",
      categories: "分类",
      reports: "报表",
    },
    overview: {
      title: "财务概览",
      subtitle: "一目了然掌握您的财务状况",
      quickActions: "快速操作",
      viewAllAccounts: "查看所有账户",
      manageYourAccounts: "管理您的银行账户和信用卡",
      viewAllTransactions: "查看所有交易",
      browseRecentTransactions: "浏览和管理最近的交易记录",
      importStatement: "导入银行对账单",
      importBankStatement: "从银行CSV或OFX文件导入交易",
      viewReports: "查看报表",
      analyzeYourFinance: "分析您的收入、支出和趋势",
    },
    summary: {
      totalBalance: "总余额",
      totalIncome: "总收入",
      totalExpense: "总支出",
      netCashFlow: "净现金流",
    },
    actions: {
      addTransaction: "添加交易",
      addAccount: "添加账户",
    },
    accounts: {
      title: "账户管理",
      addAccount: "添加账户",
      table: {
        name: "账户名称",
        type: "类型",
        currency: "货币",
        balance: "余额",
      },
      types: {
        checking: "支票账户",
        savings: "储蓄账户",
        credit: "信用卡",
        investment: "投资账户",
      },
      errors: {
        loadFailed: "加载账户失败",
      },
    },
    transactions: {
      title: "交易记录",
      addTransaction: "添加交易",
    },
    import: {
      title: "导入银行对账单",
    },
    categories: {
      title: "分类管理",
    },
    reports: {
      title: "财务报表",
    },
    filters: {
      dateRange: "选择日期范围",
      allAccounts: "所有账户",
      allTypes: "所有类型",
      income: "收入",
      expense: "支出",
      transfer: "转账",
      adjustment: "调整",
    },
    table: {
      date: "日期",
      account: "账户",
      type: "类型",
      description: "描述",
      category: "分类",
      amount: "金额",
    },
    empty: {
      noTransactions: "暂无交易记录",
      addFirstTransaction: "添加首笔交易",
    },
  },

  // 设置
  settings: {
    index: {
      description: "管理您的应用程序偏好和配置",
    },
    general: {
      title: "通用",
      description: "语言、货币、日期格式和其他偏好设置",
    },
    profile: {
      title: "个人资料",
      description: "您的个人信息和公司详细信息",
    },
    invoice: {
      title: "发票",
      description: "发票默认设置、条款和抬头信息",
    },
    email: {
      title: "邮件",
      description: "邮件提供商设置和模板",
    },
    finance: {
      title: "财务",
      description: "财务账户和交易设置",
    },
  },
};
