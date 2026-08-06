package githublib

// RecommendedRepo is a curated sysadmin GitHub repository.
// The URL is never sent to the frontend (json:"-") — installs are resolved
// server-side by repo name so end-users never see GitHub links.
type RecommendedRepo struct {
	Name     string `json:"name"`
	URL      string `json:"-"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
}

// RecommendedRepos is the curated list shown in the GitHub Library.
var RecommendedRepos = []*RecommendedRepo{
	// Active Directory
	{Name: "Adalanche", URL: "https://github.com/lkarlslund/Adalanche", Category: "Active Directory", Summary: "Active Directory ACL analysis"},
	{Name: "AppGroup", URL: "https://github.com/iandiv/AppGroup", Category: "Active Directory", Summary: "AD group management utility"},
	{Name: "DSInternals", URL: "https://github.com/MichaelGrafnetter/DSInternals", Category: "Active Directory", Summary: "Active Directory diagnostics"},
	{Name: "Graphpython", URL: "https://github.com/mlcsec/Graphpython", Category: "Active Directory", Summary: "Microsoft Graph API Python toolkit"},
	{Name: "M365Corner-Scripts", URL: "https://github.com/m365corner/M365Corner-Scripts", Category: "Active Directory", Summary: "Entra apps auditing and governance scripts"},

	// AI & Machine Learning
	{Name: "500-AI-Machine-learning-Deep-learning-Computer-vision-NLP-Projects-with-code", URL: "https://github.com/ashishpatel26/500-AI-Machine-learning-Deep-learning-Computer-vision-NLP-Projects-with-code", Category: "AI & Machine Learning", Summary: "500+ AI/ML projects with code"},
	{Name: "AI-For-Beginners", URL: "https://github.com/microsoft/AI-For-Beginners", Category: "AI & Machine Learning", Summary: "Microsoft AI curriculum for beginners"},
	{Name: "AutoGPT", URL: "https://github.com/Significant-Gravitas/AutoGPT", Category: "AI & Machine Learning", Summary: "Autonomous GPT-4 agent"},
	{Name: "awesome-ai-agents", URL: "https://github.com/e2b-dev/awesome-ai-agents", Category: "AI & Machine Learning", Summary: "Curated list of AI agents"},
	{Name: "chatbox", URL: "https://github.com/chatboxai/chatbox", Category: "AI & Machine Learning", Summary: "Desktop AI client for multiple LLMs"},
	{Name: "claude-mem", URL: "https://github.com/thedotmack/claude-mem", Category: "AI & Machine Learning", Summary: "Memory management for Claude AI"},
	{Name: "context-mode", URL: "https://github.com/mksglu/context-mode", Category: "AI & Machine Learning", Summary: "Context mode for AI coding"},
	{Name: "dify", URL: "https://github.com/langgenius/dify", Category: "AI & Machine Learning", Summary: "LLM app development platform"},
	{Name: "Flowise", URL: "https://github.com/FlowiseAI/Flowise", Category: "AI & Machine Learning", Summary: "Drag and drop LLM workflow builder"},
	{Name: "hermes-agent", URL: "https://github.com/NousResearch/hermes-agent", Category: "AI & Machine Learning", Summary: "Hermes AI agent framework"},
	{Name: "kilocode", URL: "https://github.com/Kilo-Org/kilocode", Category: "AI & Machine Learning", Summary: "AI coding assistant for VS Code"},
	{Name: "langchain", URL: "https://github.com/langchain-ai/langchain", Category: "AI & Machine Learning", Summary: "LLM application framework"},
	{Name: "learn-claude-code", URL: "https://github.com/shareAI-lab/learn-claude-code", Category: "AI & Machine Learning", Summary: "Claude coding tutorials"},
	{Name: "llama.cpp", URL: "https://github.com/ggml-org/llama.cpp", Category: "AI & Machine Learning", Summary: "Run LLMs locally in C/C++"},
	{Name: "LLMs-from-scratch", URL: "https://github.com/rasbt/LLMs-from-scratch", Category: "AI & Machine Learning", Summary: "Build LLMs from scratch tutorials"},
	{Name: "minimind", URL: "https://github.com/jingyaogong/minimind", Category: "AI & Machine Learning", Summary: "Mini mind LLM from scratch"},
	{Name: "NextChat", URL: "https://github.com/ChatGPTNextWeb/NextChat", Category: "AI & Machine Learning", Summary: "ChatGPT-like client"},
	{Name: "n8n", URL: "https://github.com/n8n-io/n8n", Category: "AI & Machine Learning", Summary: "Workflow automation platform"},
	{Name: "ollama", URL: "https://github.com/ollama/ollama", Category: "AI & Machine Learning", Summary: "Run LLMs locally"},
	{Name: "OpenHands", URL: "https://github.com/OpenHands/OpenHands", Category: "AI & Machine Learning", Summary: "AI-powered coding agent"},
	{Name: "prompts.chat", URL: "https://github.com/f/prompts.chat", Category: "AI & Machine Learning", Summary: "Prompt engineering playground"},
	{Name: "spaCy", URL: "https://github.com/explosion/spaCy", Category: "AI & Machine Learning", Summary: "Industrial NLP library"},
	{Name: "system-prompts-and-models-of-ai-tools", URL: "https://github.com/x1xhlol/system-prompts-and-models-of-ai-tools", Category: "AI & Machine Learning", Summary: "System prompts for AI tools"},
	{Name: "SWE-agent", URL: "https://github.com/SWE-agent/SWE-agent", Category: "AI & Machine Learning", Summary: "AI software engineering agent"},
	{Name: "transformers", URL: "https://github.com/huggingface/transformers", Category: "AI & Machine Learning", Summary: "HuggingFace transformer models"},
	{Name: "vllm", URL: "https://github.com/vllm-project/vllm", Category: "AI & Machine Learning", Summary: "High-throughput LLM inference"},

	// Azure & Cloud
	{Name: "azure-templates", URL: "https://github.com/fortinet/azure-templates", Category: "Azure & Cloud", Summary: "Fortinet Azure templates"},
	{Name: "cloud-custodian", URL: "https://github.com/cloud-custodian/cloud-custodian", Category: "Azure & Cloud", Summary: "Cloud governance-as-code"},
	{Name: "cloudformation-guard", URL: "https://github.com/aws-cloudformation/cloudformation-guard", Category: "Azure & Cloud", Summary: "CloudFormation compliance rules"},

	// Databases
	{Name: "beekeeper-studio", URL: "https://github.com/beekeeper-studio/beekeeper-studio", Category: "Databases", Summary: "Cross-platform database GUI"},
	{Name: "dbgate", URL: "https://github.com/dbgate/dbgate", Category: "Databases", Summary: "Database management GUI"},
	{Name: "fluentmigrator", URL: "https://github.com/fluentmigrator/fluentmigrator", Category: "Databases", Summary: "Database migration framework"},
	{Name: "immudb", URL: "https://github.com/codenotary/immudb", Category: "Databases", Summary: "Immutable database for audit"},
	{Name: "querybuilder", URL: "https://github.com/sqlkata/querybuilder", Category: "Databases", Summary: "SQL query builder for Go"},

	// DevOps & Labs
	{Name: "ansible", URL: "https://github.com/ansible/ansible", Category: "DevOps & Labs", Summary: "IT automation platform"},
	{Name: "ansible-collection-hardening", URL: "https://github.com/dev-sec/ansible-collection-hardening", Category: "DevOps & Labs", Summary: "Security hardening for Ansible"},
	{Name: "ansible-hyperv", URL: "https://github.com/tsailiming/ansible-hyperv", Category: "DevOps & Labs", Summary: "Ansible Hyper-V collection"},
	{Name: "AutomatedBadLab", URL: "https://github.com/spyr0-sec/AutomatedBadLab", Category: "DevOps & Labs", Summary: "Automated lab builder"},
	{Name: "AutomatedLab", URL: "https://github.com/AutomatedLab/AutomatedLab", Category: "DevOps & Labs", Summary: "Lab automation for Hyper-V and Azure"},
	{Name: "Datto_Powershell", URL: "https://github.com/westgate-computers/Datto_Powershell", Category: "DevOps & Labs", Summary: "Datto RMM PowerShell scripts"},
	{Name: "DetectionLab", URL: "https://github.com/clong/DetectionLab", Category: "DevOps & Labs", Summary: "Security detection lab"},
	{Name: "fortinet-zabbix", URL: "https://github.com/mbdraks/fortinet-zabbix", Category: "DevOps & Labs", Summary: "Fortinet Zabbix integration"},
	{Name: "kpt", URL: "https://github.com/kptdev/kpt", Category: "DevOps & Labs", Summary: "Kubernetes config management"},
	{Name: "kyverno", URL: "https://github.com/kyverno/kyverno", Category: "DevOps & Labs", Summary: "Kubernetes policy engine"},
	{Name: "lisa", URL: "https://github.com/microsoft/lisa", Category: "DevOps & Labs", Summary: "Linux test automation from Microsoft"},
	{Name: "MacHyperVSupport", URL: "https://github.com/acidanthera/MacHyperVSupport", Category: "DevOps & Labs", Summary: "Hyper-V support for macOS"},
	{Name: "OSX-Hyper-V", URL: "https://github.com/Qonfused/OSX-Hyper-V", Category: "DevOps & Labs", Summary: "macOS on Hyper-V guide"},
	{Name: "packer-plugin-hyperv", URL: "https://github.com/hashicorp/packer-plugin-hyperv", Category: "DevOps & Labs", Summary: "Packer plugin for Hyper-V"},
	{Name: "windows-vagrant", URL: "https://github.com/rgl/windows-vagrant", Category: "DevOps & Labs", Summary: "Windows VMs with Vagrant"},
	{Name: "VMAware", URL: "https://github.com/NotRequiem/VMAware", Category: "DevOps & Labs", Summary: "VM detection utility"},

	// GRC & Compliance
	{Name: "ciso-assistant-community", URL: "https://github.com/intuitem/ciso-assistant-community", Category: "GRC & Compliance", Summary: "CISO risk management tool"},
	{Name: "claude-grc-engineering", URL: "https://github.com/GRCEngClub/claude-grc-engineering", Category: "GRC & Compliance", Summary: "GRC engineering with Claude AI"},
	{Name: "comply", URL: "https://github.com/strongdm/comply", Category: "GRC & Compliance", Summary: "IT compliance framework"},
	{Name: "ECC", URL: "https://github.com/affaan-m/ECC", Category: "GRC & Compliance", Summary: "Endpoint compliance checker"},
	{Name: "pacbot", URL: "https://github.com/tmobile/pacbot", Category: "GRC & Compliance", Summary: "Cloud compliance scanner from T-Mobile"},
	{Name: "probo", URL: "https://github.com/getprobo/probo", Category: "GRC & Compliance", Summary: "Open source GRC platform"},
	{Name: "ScubaGear", URL: "https://github.com/cisagov/ScubaGear", Category: "GRC & Compliance", Summary: "M365 security assessment from CISA"},
	{Name: "verifywise", URL: "https://github.com/verifywise-ai/verifywise", Category: "GRC & Compliance", Summary: "AI-powered compliance verification"},

	// Microsoft 365
	{Name: "m365apps-deploy", URL: "https://github.com/haakonwibe/m365apps-deploy", Category: "Microsoft 365", Summary: "Deploy M365 Apps, Visio, Project via Intune Win32"},
	{Name: "SharePoint-Smart-Copy", URL: "https://github.com/sregan1/SharePoint-Smart-Copy", Category: "Microsoft 365", Summary: "Copy files between SharePoint site collections"},
	{Name: "SharePoint-Smart-Org-Chart", URL: "https://github.com/sregan1/SharePoint-Smart-Org-Chart", Category: "Microsoft 365", Summary: "SPFx org chart web part with Graph"},
	{Name: "m365corner-reporting-tool", URL: "https://github.com/m365corner/m365corner-reporting-tool-community-edition", Category: "Microsoft 365", Summary: "M365 reporting for users, groups, Teams"},
	{Name: "aegis", URL: "https://github.com/anthonyonazure/aegis", Category: "Microsoft 365", Summary: "AI-powered M365 governance for MSPs"},
	{Name: "meeting-room-manager", URL: "https://github.com/AliMalik5615/meeting-room-manager", Category: "Microsoft 365", Summary: "Exchange meeting room provisioning"},
	{Name: "Online-Archive-Activation", URL: "https://github.com/broccoliandpepper/Online-Archive-Activation", Category: "Microsoft 365", Summary: "Enable Exchange Online archives for a group"},
	{Name: "M365wizard", URL: "https://github.com/DwayneSelsig/M365wizard", Category: "Microsoft 365", Summary: "M365 decision guidance for IT teams"},
	{Name: "lily-pwsh", URL: "https://github.com/revpixel/lily-pwsh", Category: "Microsoft 365", Summary: "Zero-drift PowerShell container for M365"},
	{Name: "aether365", URL: "https://github.com/aether365-io/cli", Category: "Microsoft 365", Summary: "M365 compliance scanner (CIS/NIS2)"},
	{Name: "Siemserva", URL: "https://github.com/Senserva/Siemserva-Releases", Category: "Microsoft 365", Summary: "Microsoft-focused security management"},

	// AI & Machine Learning
	{Name: "ai-agents-for-beginners", URL: "https://github.com/microsoft/ai-agents-for-beginners", Category: "AI & Machine Learning", Summary: "Microsoft course: build AI agents (18 lessons)"},

	// Intune & Endpoint
	{Name: "AdminToolbox", URL: "https://github.com/TheTaylorLee/AdminToolbox", Category: "Intune & Endpoint", Summary: "IT admin PowerShell toolkit"},
	{Name: "AutotaskAPI", URL: "https://github.com/KelvinTegelaar/AutotaskAPI", Category: "Intune & Endpoint", Summary: "Autotask API PowerShell module"},
	{Name: "Celerium.DattoRMM", URL: "https://github.com/Celerium/Celerium.DattoRMM", Category: "Intune & Endpoint", Summary: "Datto RMM PowerShell module"},
	{Name: "CustomComplianceScripts", URL: "https://github.com/JayRHa/CustomComplianceScripts", Category: "Intune & Endpoint", Summary: "Intune custom compliance scripts"},
	{Name: "Datto-PowerShellWrapper", URL: "https://github.com/Celerium/Datto-PowerShellWrapper", Category: "Intune & Endpoint", Summary: "Datto PowerShell API wrapper"},
	{Name: "Datto-RMM-UI-Enhancer", URL: "https://github.com/mdjx/Datto-RMM-UI-Enhancer", Category: "Intune & Endpoint", Summary: "Datto RMM UI enhancement"},
	{Name: "EndpointAnalyticsRemediationScripts", URL: "https://github.com/JayRHa/EndpointAnalyticsRemediationScripts", Category: "Intune & Endpoint", Summary: "Intune analytics remediation"},
	{Name: "intune-patching-os-compliance-dashboard", URL: "https://github.com/greebo-labs/intune-patching-os-compliance-dashboard", Category: "Intune & Endpoint", Summary: "Intune patching dashboard"},
	{Name: "IntuneCommander", URL: "https://github.com/adamgell/IntuneCommander", Category: "Intune & Endpoint", Summary: "Intune management CLI tool"},
	{Name: "IntuneCustomCompliance", URL: "https://github.com/alexverboon/IntuneCustomCompliance", Category: "Intune & Endpoint", Summary: "Intune custom compliance rules"},
	{Name: "IntuneDeviceTroubleshooter", URL: "https://github.com/JayRHa/IntuneDeviceTroubleshooter", Category: "Intune & Endpoint", Summary: "Intune device troubleshooting"},
	{Name: "IntuneScripts", URL: "https://github.com/JayRHa/IntuneScripts", Category: "Intune & Endpoint", Summary: "Intune management scripts"},
	{Name: "MEM.Zone", URL: "https://github.com/MEM-Zone/MEM.Zone", Category: "Intune & Endpoint", Summary: "Microsoft Endpoint Manager tools"},
	{Name: "ModernWorkplaceConcierge", URL: "https://github.com/nicolonsky/ModernWorkplaceConcierge", Category: "Intune & Endpoint", Summary: "Modern workplace management"},
	{Name: "rmm-scripts", URL: "https://github.com/limehawk/rmm-scripts", Category: "Intune & Endpoint", Summary: "RMM automation scripts"},
	{Name: "PSAppDeployToolkit", URL: "https://github.com/PSAppDeployToolkit/PSAppDeployToolkit", Category: "Intune & Endpoint", Summary: "Application deployment toolkit"},

	// Microsoft 365
	{Name: "cli-microsoft365", URL: "https://github.com/pnp/cli-microsoft365", Category: "Microsoft 365", Summary: "CLI for managing M365 tenant"},
	{Name: "content", URL: "https://github.com/ComplianceAsCode/content", Category: "Microsoft 365", Summary: "Compliance as Code content"},
	{Name: "Connect-MS365", URL: "https://github.com/blindzero/Connect-MS365", Category: "Microsoft 365", Summary: "M365 connection helper"},
	{Name: "dev-proxy", URL: "https://github.com/dotnet/dev-proxy", Category: "Microsoft 365", Summary: "Microsoft 365 API proxy for dev"},
	{Name: "Export-RecipientPermissions", URL: "https://github.com/GruberMarkus/Export-RecipientPermissions", Category: "Microsoft 365", Summary: "Export Exchange recipient permissions"},
	{Name: "List-Formatting", URL: "https://github.com/pnp/List-Formatting", Category: "Microsoft 365", Summary: "SharePoint list formatting samples"},
	{Name: "MAAD-AF", URL: "https://github.com/vectra-ai-research/MAAD-AF", Category: "Microsoft 365", Summary: "Microsoft Attack and Defense Framework"},
	{Name: "M365-Assess", URL: "https://github.com/Galvnyz/M365-Assess", Category: "Microsoft 365", Summary: "M365 tenant assessment"},
	{Name: "M365-Assessment-Toolkit", URL: "https://github.com/malcolmmcdonald1982/M365-Assessment-Toolkit", Category: "Microsoft 365", Summary: "M365 assessment toolkit"},
	{Name: "maester", URL: "https://github.com/maester365/maester", Category: "Microsoft 365", Summary: "M365 Entra ID testing framework"},
	{Name: "Microsoft365DSC", URL: "https://github.com/microsoft/Microsoft365DSC", Category: "Microsoft 365", Summary: "M365 Desired State Configuration"},
	{Name: "Microsoft365R", URL: "https://github.com/Azure/Microsoft365R", Category: "Microsoft 365", Summary: "Microsoft 365 PowerShell module"},
	{Name: "Microsoft_Cloud_Security", URL: "https://github.com/tomwechsler/Microsoft_Cloud_Security", Category: "Microsoft 365", Summary: "Azure cloud security scripts"},
	{Name: "mggraph-intune-samples", URL: "https://github.com/microsoft/mggraph-intune-samples", Category: "Microsoft 365", Summary: "MS Graph Intune samples"},
	{Name: "microsoft-graph-toolkit", URL: "https://github.com/microsoftgraph/microsoft-graph-toolkit", Category: "Microsoft 365", Summary: "Microsoft Graph Toolkit components"},
	{Name: "Microsoft-Integration-and-Azure-Stencils-Pack-for-Visio", URL: "https://github.com/sandroasp/Microsoft-Integration-and-Azure-Stencils-Pack-for-Visio", Category: "Microsoft 365", Summary: "Visio stencils for M365 and Azure"},
	{Name: "MOFA", URL: "https://github.com/cocopuff2u/MOFA", Category: "Microsoft 365", Summary: "Microsoft 365 automation framework"},
	{Name: "Office365", URL: "https://github.com/directorcia/Office365", Category: "Microsoft 365", Summary: "Office 365 PowerShell scripts"},
	{Name: "office365dev", URL: "https://github.com/chenxizhang/office365dev", Category: "Microsoft 365", Summary: "Office 365 developer tools"},
	{Name: "Office365DnsChecker", URL: "https://github.com/rhymeswithmogul/Office365DnsChecker", Category: "Microsoft 365", Summary: "Office 365 DNS records checker"},
	{Name: "O365scripts", URL: "https://github.com/O365scripts/O365scripts", Category: "Microsoft 365", Summary: "Office 365 PowerShell scripts"},
	{Name: "Office-Tool", URL: "https://github.com/YerongAI/Office-Tool", Category: "Microsoft 365", Summary: "Microsoft Office installer tool"},
	{Name: "office365-rest-python-client", URL: "https://github.com/vgrem/office365-rest-python-client", Category: "Microsoft 365", Summary: "Python client for Office 365 REST API"},
	{Name: "onedrive", URL: "https://github.com/abraunegg/onedrive", Category: "Microsoft 365", Summary: "OneDrive Linux client"},
	{Name: "OneManager-php", URL: "https://github.com/qkqpttgf/OneManager-php", Category: "Microsoft 365", Summary: "OneDrive directory manager"},
	{Name: "PnP", URL: "https://github.com/pnp/PnP", Category: "Microsoft 365", Summary: "Patterns and Practices for M365"},
	{Name: "PnP-PowerShell", URL: "https://github.com/pnp/PnP-PowerShell", Category: "Microsoft 365", Summary: "PnP PowerShell cmdlets for M365"},
	{Name: "python-o365", URL: "https://github.com/O365/python-o365", Category: "Microsoft 365", Summary: "Python library for Office 365"},
	{Name: "PS-ActiveDirectory-AzureAD-O365", URL: "https://github.com/adrianlois/PS-ActiveDirectory-AzureAD-O365", Category: "Microsoft 365", Summary: "PS scripts for AD, Azure AD, M365"},
	{Name: "PSTeams", URL: "https://github.com/EvotecIT/PSTeams", Category: "Microsoft 365", Summary: "Microsoft Teams PowerShell module"},
	{Name: "script-samples", URL: "https://github.com/pnp/script-samples", Category: "Microsoft 365", Summary: "PnP script samples for M365"},
	{Name: "sp-dev-docs", URL: "https://github.com/SharePoint/sp-dev-docs", Category: "Microsoft 365", Summary: "SharePoint developer documentation"},
	{Name: "TrustM365", URL: "https://github.com/AntoPorter/TrustM365", Category: "Microsoft 365", Summary: "Microsoft 365 trust assessment"},

	// Microsoft Defender
	{Name: "Defender-Advanced-Hunting-Queries", URL: "https://github.com/francoisfried/Defender-Advanced-Hunting-Queries", Category: "Microsoft Defender", Summary: "Advanced hunting queries for Defender"},
	{Name: "defender-asr-admx", URL: "https://github.com/MichaelGrafnetter/defender-asr-admx", Category: "Microsoft Defender", Summary: "Defender ASR rules ADMX template"},

	// Monitoring & SIEM
	{Name: "grafana", URL: "https://github.com/grafana/grafana", Category: "Monitoring & SIEM", Summary: "Observability dashboards"},
	{Name: "MISP", URL: "https://github.com/MISP/MISP", Category: "Monitoring & SIEM", Summary: "Threat intelligence platform"},
	{Name: "orca", URL: "https://github.com/cammurray/orca", Category: "Monitoring & SIEM", Summary: "M365 security posture tool"},
	{Name: "prometheus", URL: "https://github.com/prometheus/prometheus", Category: "Monitoring & SIEM", Summary: "Monitoring and alerting toolkit"},
	{Name: "prometheus-msteams", URL: "https://github.com/prometheus-msteams/prometheus-msteams", Category: "Monitoring & SIEM", Summary: "Prometheus alerts to MS Teams"},
	{Name: "wazuh", URL: "https://github.com/wazuh/wazuh", Category: "Monitoring & SIEM", Summary: "Open source security monitoring"},

	// Network & Infrastructure
	{Name: "fortinet-azure-templates", URL: "https://github.com/fortinet/azure-templates", Category: "Network & Infrastructure", Summary: "Fortinet Azure deployment templates"},
	{Name: "NET-AUTOMATE", URL: "https://github.com/kebaldwi/NET-AUTOMATE", Category: "Network & Infrastructure", Summary: "Network automation scripts"},
	{Name: "netbox", URL: "https://github.com/netbox-community/netbox", Category: "Network & Infrastructure", Summary: "IP address management tool"},
	{Name: "Netshot", URL: "https://github.com/netfishers-onl/Netshot", Category: "Network & Infrastructure", Summary: "Network device configuration management"},
	{Name: "network-mcp-docker-suite", URL: "https://github.com/pamosima/network-mcp-docker-suite", Category: "Network & Infrastructure", Summary: "Network MCP Docker tools"},
	{Name: "PSCiscoMeraki", URL: "https://github.com/BanterBoy/PSCiscoMeraki", Category: "Network & Infrastructure", Summary: "Cisco Meraki PowerShell module"},
	{Name: "terraform-provider-meraki", URL: "https://github.com/CiscoDevNet/terraform-provider-meraki", Category: "Network & Infrastructure", Summary: "Terraform provider for Meraki"},

	// Security & Monitoring
	{Name: "fixinventory", URL: "https://github.com/someengineering/fixinventory", Category: "Security & Monitoring", Summary: "Cloud infrastructure inventory"},
	{Name: "trivy", URL: "https://github.com/aquasecurity/trivy", Category: "Security & Monitoring", Summary: "Vulnerability scanner for containers and OS"},

	// Utilities & Tools
	{Name: "awesome-sysadmin", URL: "https://github.com/n1trux/awesome-sysadmin", Category: "Utilities & Tools", Summary: "Curated list of sysadmin tools"},
	{Name: "bitwarden", URL: "https://github.com/bitwarden/server", Category: "Utilities & Tools", Summary: "Open source password manager"},
	{Name: "botkit", URL: "https://github.com/howdyai/botkit", Category: "Utilities & Tools", Summary: "Bot framework toolkit"},
	{Name: "community-edition", URL: "https://github.com/ramboxapp/community-edition", Category: "Utilities & Tools", Summary: "Cross-platform app launcher"},
	{Name: "cyberduck", URL: "https://github.com/iterate-ch/cyberduck", Category: "Utilities & Tools", Summary: "Cloud storage browser"},
	{Name: "drawdb", URL: "https://github.com/drawdb-io/drawdb", Category: "Utilities & Tools", Summary: "Database diagram editor"},
	{Name: "glpi", URL: "https://github.com/glpi-project/glpi", Category: "Utilities & Tools", Summary: "IT asset management"},
	{Name: "infisical", URL: "https://github.com/Infisical/infisical", Category: "Utilities & Tools", Summary: "Secret management platform"},
	{Name: "kong", URL: "https://github.com/Kong/kong", Category: "Utilities & Tools", Summary: "API gateway"},
	{Name: "matterbridge", URL: "https://github.com/42wim/matterbridge", Category: "Utilities & Tools", Summary: "Multi-protocol chat bridge"},
	{Name: "noid-privacy", URL: "https://github.com/NexusOne23/noid-privacy", Category: "Utilities & Tools", Summary: "Privacy-first identity tool"},
	{Name: "notify", URL: "https://github.com/nikoksr/notify", Category: "Utilities & Tools", Summary: "Multi-channel notification tool"},
	{Name: "Vane", URL: "https://github.com/ItzCrazyKns/Vane", Category: "Utilities & Tools", Summary: "Web search aggregator"},

	// Windows Administration
	{Name: "choco", URL: "https://github.com/chocolatey/choco", Category: "Windows Administration", Summary: "Package manager for Windows"},
	{Name: "cmtraceopen", URL: "https://github.com/adamgell/cmtraceopen", Category: "Windows Administration", Summary: "CMTrace log viewer launcher"},
	{Name: "Debloat-Windows-10", URL: "https://github.com/W4RH4WK/Debloat-Windows-10", Category: "Windows Administration", Summary: "Debloat Windows 10"},
	{Name: "DeskHider", URL: "https://github.com/iandiv/DeskHider", Category: "Windows Administration", Summary: "Hide desktop icons quickly"},
	{Name: "DevToys", URL: "https://github.com/DevToys-app/DevToys", Category: "Windows Administration", Summary: "Developer utility toolbox for Windows"},
	{Name: "DisplayProfileManager", URL: "https://github.com/zac15987/DisplayProfileManager", Category: "Windows Administration", Summary: "Display profile switcher"},
	{Name: "DidMySettingsChange", URL: "https://github.com/nolesapex/DidMySettingsChange", Category: "Windows Administration", Summary: "Detect Windows settings changes"},
	{Name: "Files", URL: "https://github.com/files-community/Files", Category: "Windows Administration", Summary: "Modern file manager for Windows"},
	{Name: "Folcolor", URL: "https://github.com/kweatherman/Folcolor", Category: "Windows Administration", Summary: "Folder color customization"},
	{Name: "hyper-reV", URL: "https://github.com/noahware/hyper-reV", Category: "Windows Administration", Summary: "Hyper-V and Windows optimization"},
	{Name: "InstallerClean", URL: "https://github.com/no-faff/InstallerClean", Category: "Windows Administration", Summary: "Clean up installer leftovers"},
	{Name: "lively", URL: "https://github.com/rocksdanister/lively", Category: "Windows Administration", Summary: "Animated desktop wallpapers"},
	{Name: "NanaZip", URL: "https://github.com/M2Team/NanaZip", Category: "Windows Administration", Summary: "Modern archive manager for Windows"},
	{Name: "PowerShell", URL: "https://github.com/fleschutz/PowerShell", Category: "Windows Administration", Summary: "1000+ essential PowerShell scripts"},
	{Name: "PowerShell-collection", URL: "https://github.com/jhochwald/PowerShell-collection", Category: "Windows Administration", Summary: "Curated PowerShell script collection"},
	{Name: "PowerToys", URL: "https://github.com/microsoft/PowerToys", Category: "Windows Administration", Summary: "Windows power user tools"},
	{Name: "PSScriptAnalyzer", URL: "https://github.com/PowerShell/PSScriptAnalyzer", Category: "Windows Administration", Summary: "PowerShell code linter and checker"},
	{Name: "Public-Scripts", URL: "https://github.com/Mike-Crowley/Public-Scripts", Category: "Windows Administration", Summary: "Public IT admin scripts"},
	{Name: "Seelen-UI", URL: "https://github.com/eythaann/Seelen-UI", Category: "Windows Administration", Summary: "Desktop customization tool"},
	{Name: "UniGetUI", URL: "https://github.com/Devolutions/UniGetUI", Category: "Windows Administration", Summary: "GUI for package managers"},
	{Name: "W1X-Debloat", URL: "https://github.com/AdminVin/W1X-Debloat", Category: "Windows Administration", Summary: "Windows 10/11 debloat tool"},
	{Name: "Win-Debloat-Tools", URL: "https://github.com/LeDragoX/Win-Debloat-Tools", Category: "Windows Administration", Summary: "Windows debloat and optimization"},
	{Name: "Win11Debloat", URL: "https://github.com/Raphire/Win11Debloat", Category: "Windows Administration", Summary: "Remove bloatware from Windows 11"},
	{Name: "Win11Tweaker", URL: "https://github.com/iandiv/Win11Tweaker", Category: "Windows Administration", Summary: "Windows 11 tweaking tool"},
	{Name: "Windows-11-Guide", URL: "https://github.com/mikeroyal/Windows-11-Guide", Category: "Windows Administration", Summary: "Windows 11 setup guide"},
	{Name: "Windows_Toolz", URL: "https://github.com/Exploitacious/Windows_Toolz", Category: "Windows Administration", Summary: "Windows utility scripts"},
	{Name: "WindowsToolbox", URL: "https://github.com/WinTweakers/WindowsToolbox", Category: "Windows Administration", Summary: "Windows maintenance toolbox"},
	{Name: "Winhance", URL: "https://github.com/memstechtips/Winhance", Category: "Windows Administration", Summary: "Windows optimization and enhancement"},
	{Name: "winget-cli", URL: "https://github.com/microsoft/winget-cli", Category: "Windows Administration", Summary: "Windows Package Manager CLI"},
}

func init() {
	for _, r := range RecommendedRepos {
		RegisterURL(r.Name, r.URL)
	}
}

// GetRecommendedRepos returns the curated sysadmin repo list.
func GetRecommendedRepos() []*RecommendedRepo {
	return RecommendedRepos
}
