"use client";

import React, { Suspense, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { jwtDecode } from "jwt-decode";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Team } from "@/components/key_team_helpers/key_list";
import Navbar from "@/components/navbar";
import UserDashboard from "@/components/user_dashboard";
import ModelDashboard from "@/components/model_dashboard";
import ViewUserDashboard from "@/components/view_users";
import Teams from "@/components/teams";
import Organizations from "@/components/organizations";
import { fetchOrganizations } from "@/components/organizations";
import AdminPanel from "@/components/admins";
import Settings from "@/components/settings";
import GeneralSettings from "@/components/general_settings";
import PassThroughSettings from "@/components/pass_through_settings";
import BudgetPanel from "@/components/budgets/budget_panel";
import SpendLogsTable from "@/components/view_logs";
import ModelHubTable from "@/components/model_hub_table";
import NewUsagePage from "@/components/new_usage";
import APIRef from "@/components/api_ref";
import ChatUI from "@/components/chat_ui";
import Sidebar from "@/components/leftnav";
import Usage from "@/components/usage";
import CacheDashboard from "@/components/cache_dashboard";
import {
  getUiConfig,
  proxyBaseUrl,
  setGlobalLitellmHeaderName,
} from "@/components/networking";
import { Organization } from "@/components/networking";
import GuardrailsPanel from "@/components/guardrails";
import TransformRequestPanel from "@/components/transform_request";
import { fetchUserModels } from "@/components/create_key_button";
import { fetchTeams } from "@/components/common_components/fetch_teams";
import { MCPServers } from "@/components/mcp_tools";
import TagManagement from "@/components/tag_management";
import VectorStoreManagement from "@/components/vector_store_management";
import { UiLoadingSpinner } from "@/components/ui/ui-loading-spinner";
import { cx } from "@/lib/cva.config";

function getCookie(name: string) {
  const cookieValue = document.cookie
    .split("; ")
    .find((row) => row.startsWith(name + "="));
  return cookieValue ? cookieValue.split("=")[1] : null;
}

function formatUserRole(userRole: string) {
  if (!userRole) {
    return "Undefined Role";
  }
  switch (userRole.toLowerCase()) {
    case "app_owner":
      return "App Owner";
    case "demo_app_owner":
      return "App Owner";
    case "app_admin":
      return "Admin";
    case "proxy_admin":
      return "Admin";
    case "proxy_admin_viewer":
      return "Admin Viewer";
    case "org_admin":
      return "Org Admin";
    case "internal_user":
      return "Internal User";
    case "internal_user_viewer":
    case "internal_viewer": // TODO:remove if deprecated
      return "Internal Viewer";
    case "app_user":
      return "App User";
    default:
      return "Unknown Role";
  }
}

interface ProxySettings {
  PROXY_BASE_URL: string;
  PROXY_LOGOUT_URL: string;
}

const queryClient = new QueryClient();

function LoadingScreen() {
  return (
    <div className={cx("h-screen", "flex items-center justify-center gap-4")}>
      <div className="text-lg font-medium py-2 pr-4 border-r border-r-gray-200">
        🚅 LiteLLM
      </div>

      <div className="flex items-center justify-center gap-2">
        <UiLoadingSpinner className="size-4" />
        <span className="text-gray-600 text-sm">Loading...</span>
      </div>
    </div>
  );
}

export default function CreateKeyPage() {
  const [userRole, setUserRole] = useState("");
  const [premiumUser, setPremiumUser] = useState(false);
  const [disabledPersonalKeyCreation, setDisabledPersonalKeyCreation] =
    useState(false);
  const [userEmail, setUserEmail] = useState<null | string>(null);
  const [teams, setTeams] = useState<Team[] | null>(null);
  const [keys, setKeys] = useState<null | any[]>([]);
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [userModels, setUserModels] = useState<string[]>([]);
  const [proxySettings, setProxySettings] = useState<ProxySettings>({
    PROXY_BASE_URL: "",
    PROXY_LOGOUT_URL: "",
  });

  const [showSSOBanner, setShowSSOBanner] = useState<boolean>(true);
  const searchParams = useSearchParams()!;
  const [modelData, setModelData] = useState<any>({ data: [] });
  const [token, setToken] = useState<string | null>(null);
  const [createClicked, setCreateClicked] = useState<boolean>(false);
  const [authLoading, setAuthLoading] = useState(true);
  const [userID, setUserID] = useState<string | null>(null);

  const invitation_id = searchParams.get("invitation_id");

  // Get page from URL, default to 'api-keys' if not present
  const [page, setPage] = useState(() => {
    return searchParams.get("page") || "api-keys";
  });

  // Custom setPage function that updates URL
  const updatePage = (newPage: string) => {
    // Update URL without full page reload
    const newSearchParams = new URLSearchParams(searchParams);
    newSearchParams.set("page", newPage);

    // Use Next.js router to update URL
    window.history.pushState(null, "", `?${newSearchParams.toString()}`);

    setPage(newPage);
  };

  const [accessToken, setAccessToken] = useState<string | null>(null);

  const addKey = (data: any) => {
    setKeys((prevData) => (prevData ? [...prevData, data] : [data]));
    setCreateClicked(() => !createClicked);
  };
  const redirectToLogin =
    authLoading === false && token === null && invitation_id === null;

  useEffect(() => {
    const token = getCookie("token");
    getUiConfig().then((data) => {
      // get the information for constructing the proxy base url, and then set the token and auth loading
      setToken(token);
      setAuthLoading(false);
    });
  }, []);

  useEffect(() => {
    if (redirectToLogin) {
      window.location.href = (proxyBaseUrl || "") + "/sso/key/generate";
    }
  }, [redirectToLogin]);

  useEffect(() => {
    if (!token) {
      return;
    }

    const decoded = jwtDecode(token) as { [key: string]: any };
    if (decoded) {
      // set accessToken
      setAccessToken(decoded.key);

      setDisabledPersonalKeyCreation(
        decoded.disabled_non_admin_personal_key_creation
      );

      // check if userRole is defined
      if (decoded.user_role) {
        const formattedUserRole = formatUserRole(decoded.user_role);
        setUserRole(formattedUserRole);
        if (formattedUserRole == "Admin Viewer") {
          setPage("usage");
        }
      }

      if (decoded.user_email) {
        setUserEmail(decoded.user_email);
      }

      if (decoded.login_method) {
        setShowSSOBanner(
          decoded.login_method == "username_password" ? true : false
        );
      }

      if (decoded.premium_user) {
        setPremiumUser(decoded.premium_user);
      }

      if (decoded.auth_header_name) {
        setGlobalLitellmHeaderName(decoded.auth_header_name);
      }

      if (decoded.user_id) {
        setUserID(decoded.user_id);
      }
    }
  }, [token]);

  useEffect(() => {
    if (accessToken && userID && userRole) {
      fetchUserModels(userID, userRole, accessToken, setUserModels);
    }
    if (accessToken && userID && userRole) {
      fetchTeams(accessToken, userID, userRole, null, setTeams);
    }
    if (accessToken) {
      fetchOrganizations(accessToken, setOrganizations);
    }
  }, [accessToken, userID, userRole]);

  if (authLoading || redirectToLogin) {
    return <LoadingScreen />;
  }

  return (
    <Suspense fallback={<LoadingScreen />}>
      <QueryClientProvider client={queryClient}>
        {invitation_id ? (
          <UserDashboard
            userID={userID}
            userRole={userRole}
            premiumUser={premiumUser}
            teams={teams}
            keys={keys}
            setUserRole={setUserRole}
            userEmail={userEmail}
            setUserEmail={setUserEmail}
            setTeams={setTeams}
            setKeys={setKeys}
            organizations={organizations}
            addKey={addKey}
            createClicked={createClicked}
          />
        ) : (
          <div className="flex flex-col min-h-screen">
            <Navbar
              userID={userID}
              userRole={userRole}
              premiumUser={premiumUser}
              userEmail={userEmail}
              setProxySettings={setProxySettings}
              proxySettings={proxySettings}
              accessToken={accessToken}
              isPublicPage={false}
            />
            <div className="flex flex-1 overflow-auto">
              <div className="mt-8">
                <Sidebar
                  accessToken={accessToken}
                  setPage={updatePage}
                  userRole={userRole}
                  defaultSelectedKey={page}
                />
              </div>

              {page == "api-keys" ? (
                <UserDashboard
                  userID={userID}
                  userRole={userRole}
                  premiumUser={premiumUser}
                  teams={teams}
                  keys={keys}
                  setUserRole={setUserRole}
                  userEmail={userEmail}
                  setUserEmail={setUserEmail}
                  setTeams={setTeams}
                  setKeys={setKeys}
                  organizations={organizations}
                  addKey={addKey}
                  createClicked={createClicked}
                />
              ) : page == "models" ? (
                <ModelDashboard
                  userID={userID}
                  userRole={userRole}
                  token={token}
                  keys={keys}
                  accessToken={accessToken}
                  modelData={modelData}
                  setModelData={setModelData}
                  premiumUser={premiumUser}
                  teams={teams}
                />
              ) : page == "llm-playground" ? (
                <ChatUI
                  userID={userID}
                  userRole={userRole}
                  token={token}
                  accessToken={accessToken}
                  disabledPersonalKeyCreation={disabledPersonalKeyCreation}
                />
              ) : page == "users" ? (
                <ViewUserDashboard
                  userID={userID}
                  userRole={userRole}
                  token={token}
                  keys={keys}
                  teams={teams}
                  accessToken={accessToken}
                  setKeys={setKeys}
                />
              ) : page == "teams" ? (
                <Teams
                  teams={teams}
                  setTeams={setTeams}
                  searchParams={searchParams}
                  accessToken={accessToken}
                  userID={userID}
                  userRole={userRole}
                  organizations={organizations}
                  premiumUser={premiumUser}
                />
              ) : page == "organizations" ? (
                <Organizations
                  organizations={organizations}
                  setOrganizations={setOrganizations}
                  userModels={userModels}
                  accessToken={accessToken}
                  userRole={userRole}
                  premiumUser={premiumUser}
                />
              ) : page == "admin-panel" ? (
                <AdminPanel
                  setTeams={setTeams}
                  searchParams={searchParams}
                  accessToken={accessToken}
                  userID={userID}
                  showSSOBanner={showSSOBanner}
                  premiumUser={premiumUser}
                  proxySettings={proxySettings}
                />
              ) : page == "api_ref" ? (
                <APIRef proxySettings={proxySettings} />
              ) : page == "settings" ? (
                <Settings
                  userID={userID}
                  userRole={userRole}
                  accessToken={accessToken}
                  premiumUser={premiumUser}
                />
              ) : page == "budgets" ? (
                <BudgetPanel accessToken={accessToken} />
              ) : page == "guardrails" ? (
                <GuardrailsPanel
                  accessToken={accessToken}
                  userRole={userRole}
                />
              ) : page == "transform-request" ? (
                <TransformRequestPanel accessToken={accessToken} />
              ) : page == "general-settings" ? (
                <GeneralSettings
                  userID={userID}
                  userRole={userRole}
                  accessToken={accessToken}
                  modelData={modelData}
                />
              ) : page == "model-hub-table" ? (
                <ModelHubTable
                  accessToken={accessToken}
                  publicPage={false}
                  premiumUser={premiumUser}
                  userRole={userRole}
                />
              ) : page == "caching" ? (
                <CacheDashboard
                  userID={userID}
                  userRole={userRole}
                  token={token}
                  accessToken={accessToken}
                  premiumUser={premiumUser}
                />
              ) : page == "pass-through-settings" ? (
                <PassThroughSettings
                  userID={userID}
                  userRole={userRole}
                  accessToken={accessToken}
                  modelData={modelData}
                />
              ) : page == "logs" ? (
                <SpendLogsTable
                  userID={userID}
                  userRole={userRole}
                  token={token}
                  accessToken={accessToken}
                  allTeams={(teams as Team[]) ?? []}
                  premiumUser={premiumUser}
                />
              ) : page == "mcp-servers" ? (
                <MCPServers
                  accessToken={accessToken}
                  userRole={userRole}
                  userID={userID}
                />
              ) : page == "tag-management" ? (
                <TagManagement
                  accessToken={accessToken}
                  userRole={userRole}
                  userID={userID}
                />
              ) : page == "vector-stores" ? (
                <VectorStoreManagement
                  accessToken={accessToken}
                  userRole={userRole}
                  userID={userID}
                />
              ) : page == "new_usage" ? (
                <NewUsagePage
                  userID={userID}
                  userRole={userRole}
                  accessToken={accessToken}
                  teams={(teams as Team[]) ?? []}
                  premiumUser={premiumUser}
                />
              ) : (
                <Usage
                  userID={userID}
                  userRole={userRole}
                  token={token}
                  accessToken={accessToken}
                  keys={keys}
                  premiumUser={premiumUser}
                />
              )}
            </div>
          </div>
        )}
      </QueryClientProvider>
    </Suspense>
  );
}
