import { LogoutOutlined, PlusOutlined, TeamOutlined, UnorderedListOutlined, UserOutlined } from '@ant-design/icons'
import { Layout, Menu, Space, Typography, Button } from 'antd'
import type { MenuProps } from 'antd'
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom'

import { useAuthStore } from '@/adapters/in/state/auth.store'
import styles from './AppLayout.module.css'

const { Header, Sider, Content } = Layout

function resolveSelectedKey(pathname: string): string {
  if (pathname.startsWith('/tickets/new')) return '/tickets/new'
  if (pathname.startsWith('/tickets')) return '/tickets'
  if (pathname.startsWith('/teams')) return '/teams'
  if (pathname.startsWith('/profile')) return '/profile'
  return '/tickets'
}

export function AppLayout() {
  const location = useLocation()
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)

  const selectedKey = resolveSelectedKey(location.pathname)

  const menuItems: MenuProps['items'] = [
    { key: '/tickets', icon: <UnorderedListOutlined />, label: <Link to="/tickets">Tickets</Link> },
    { key: '/tickets/new', icon: <PlusOutlined />, label: <Link to="/tickets/new">New ticket</Link> },
    { key: '/profile', icon: <UserOutlined />, label: <Link to="/profile">Profile</Link> },
  ]

  if (user?.role === 'admin') {
    menuItems.push({ key: '/teams', icon: <TeamOutlined />, label: <Link to="/teams">Teams</Link> })
  }

  return (
    <Layout className={styles.layout}>
      <Sider breakpoint="lg" collapsedWidth="0">
        <div className={styles.brand}>CoreLog</div>
        <Menu theme="dark" mode="inline" selectedKeys={[selectedKey]} items={menuItems} />
      </Sider>
      <Layout>
        <Header className={styles.header}>
          <Space>
            <Link to="/profile">
              <Typography.Text>{user?.name}</Typography.Text>
            </Link>
            <Button
              icon={<LogoutOutlined />}
              onClick={() => {
                logout()
                navigate('/login')
              }}
            >
              Log out
            </Button>
          </Space>
        </Header>
        <Content className={styles.content}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
