import React from 'react';
import { NavLink, Outlet, useLocation } from 'react-router-dom';
import { Navbar, NavbarBrand, NavbarContent, NavbarItem } from '@heroui/react';
import {
  LayoutDashboard,
  Tag,
  AppWindow,
  Cloud,
  Rocket,
} from 'lucide-react';
import ErrorMessage from './ErrorMessage';
import { useErrorStore } from '@/store/error';

const Layout: React.FC = () => {
  const { errorMessage, clearError } = useErrorStore();
  const location = useLocation();
  
  const menuItems = [
    {
      key: '/',
      icon: LayoutDashboard,
      label: '仪表板',
    },
    {
      key: '/versions',
      icon: Tag,
      label: '版本管理',
    },
    {
      key: '/applications',
      icon: AppWindow,
      label: '应用管理',
    },
    {
      key: '/environments',
      icon: Cloud,
      label: '环境管理',
    },
    {
      key: '/deployments',
      icon: Rocket,
      label: '部署管理',
    },
  ];

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path);
  };

  return (
    <div className="min-h-screen bg-background">
      <Navbar isBordered maxWidth="full" position="sticky">
        <NavbarBrand>
          <p className="font-bold text-xl">部署平台</p>
        </NavbarBrand>
        <NavbarContent className="hidden sm:flex gap-6" justify="center">
          {menuItems.map((item) => {
            const Icon = item.icon;
            return (
              <NavbarItem key={item.key} isActive={isActive(item.key)}>
                <NavLink
                  to={item.key}
                  className={({ isActive }) =>
                    `flex items-center gap-2 ${
                      isActive ? 'text-primary' : 'text-foreground'
                    }`
                  }
                >
                  <Icon size={18} />
                  <span>{item.label}</span>
                </NavLink>
              </NavbarItem>
            );
          })}
        </NavbarContent>
      </Navbar>
      <main className="container mx-auto max-w-7xl px-6 py-8">
        {errorMessage && <ErrorMessage message={errorMessage} onDismiss={clearError} />}
        <Outlet />
      </main>
    </div>
  );
};

export default Layout;
