import React from 'react';
import { useClientSettings } from '@/hooks/useClientSettings';
import { ProfileSettingsCard } from '@/components/settings/ProfileSettingsCard';
import { ApiKeysCard } from '@/components/settings/ApiKeysCard';
import { ChangePasswordCard } from '@/components/settings/ChangePasswordCard';

export const Settings: React.FC = () => {
  const {
    user,
    profileForm,
    setProfileForm,
    apiKeys,
    keyName,
    setKeyName,
    savingProfile,
    generatingKey,
    profileMessage,
    handleUpdateProfile,
    handleGenerateKey,
    handleRevokeKey,
  } = useClientSettings();

  return (
    <div className="space-y-8 animate-in fade-in-50 duration-300">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Account & Security Settings</h1>
        <p className="text-sm text-muted-foreground">
          Manage your contact information, credentials, security password, and developer API keys.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-start">
        <div className="space-y-8">
          <ProfileSettingsCard
            user={user}
            form={profileForm}
            onChange={setProfileForm}
            onSubmit={handleUpdateProfile}
            saving={savingProfile}
            message={profileMessage}
          />
          <ChangePasswordCard />
        </div>

        <ApiKeysCard
          apiKeys={apiKeys}
          keyName={keyName}
          onKeyNameChange={setKeyName}
          onGenerateKey={handleGenerateKey}
          onRevokeKey={handleRevokeKey}
          generating={generatingKey}
        />
      </div>
    </div>
  );
};

export default Settings;
