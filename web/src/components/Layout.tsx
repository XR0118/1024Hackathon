import React from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import {
  IconDashboard,
  IconTag,
  IconApps,
  IconCloud,
  IconRocket,
} from '@tabler/icons-react';

const Layout: React.FC = () => {
  const menuItems = [
    {
      key: '/',
      icon: <IconDashboard />,
      label: '仪表板',
    },
    {
      key: '/versions',
      icon: <IconTag />,
      label: '版本管理',
    },
    {
      key: '/applications',
      icon: <IconApps />,
      label: '应用管理',
    },
    {
      key: '/environments',
      icon: <IconCloud />,
      label: '环境管理',
    },
    {
      key: '/deployments',
      icon: <IconRocket />,
      label: '部署管理',
    },
  ];

  return (
    <div className="page">
      <header className="navbar navbar-expand-md d-print-none">
        <div className="container-xl">
          <h1 className="navbar-brand navbar-brand-autodark d-none-navbar-horizontal pe-0 pe-md-3">
            部署平台
          </h1>
        </div>
      </header>
      <div className="navbar-expand-md">
        <div className="collapse navbar-collapse" id="navbar-menu">
          <div className="navbar">
            <div className="container-xl">
              <ul className="navbar-nav">
                {menuItems.map((item) => (
                  <li className="nav-item" key={item.key}>
                    <NavLink className="nav-link" to={item.key} end>
                      <span className="nav-link-icon d-md-none d-lg-inline-block">
                        {item.icon}
                      </span>
                      <span className="nav-link-title">{item.label}</span>
                    </NavLink>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </div>
      <div className="page-wrapper">
        <div className="page-body">
          <div className="container-xl">
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Layout;
