import React from 'react';
import { NavLink } from 'react-router-dom';
import { navGroups } from './navConfig';

interface SidebarNavProps {
  onItemClick?: () => void;
}

export const SidebarNav: React.FC<SidebarNavProps> = ({ onItemClick }) => {
  return (
    <div className="space-y-5">
      {navGroups.map((group, gIdx) => (
        <div key={gIdx} className="space-y-1">
          {group.groupName && (
            <p className="px-3 text-[10px] font-bold uppercase tracking-wider text-muted-foreground/70 mb-1.5">
              {group.groupName}
            </p>
          )}
          {group.items.map((item) => {
            const Icon = item.icon;
            return (
              <NavLink
                key={item.path}
                to={item.path}
                end={item.path === '/'}
                onClick={onItemClick}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-3 py-2 rounded-lg text-xs font-medium transition-all ${
                    isActive
                      ? 'bg-primary text-primary-foreground shadow-sm shadow-primary/20'
                      : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                  }`
                }
              >
                <Icon className="h-4 w-4 shrink-0" />
                <span>{item.name}</span>
              </NavLink>
            );
          })}
        </div>
      ))}
    </div>
  );
};
