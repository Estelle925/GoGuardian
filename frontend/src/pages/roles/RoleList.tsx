import { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, Form, Input, Checkbox, message, Collapse, Row, Col } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, KeyOutlined } from '@ant-design/icons';
import request from '@/utils/request';
import type { ColumnsType } from 'antd/es/table';

interface Role {
  id: number;
  name: string;
  code: string;
  description?: string;
  createdAt: string;
}

interface Permission {
  id: number;
  name: string;
  enable: boolean;
  icon?: string;
  children?: Permission[];
}

const RoleList = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [permissionModalVisible, setPermissionModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const [currentRoleId, setCurrentRoleId] = useState<number | null>(null);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [permissionList, setPermissionList] = useState<Permission[]>([]);
  const [permissionTotal, setPermissionTotal] = useState(0);
  const [permissionPage, setPermissionPage] = useState(1);
  const [permissionPageSize, setPermissionPageSize] = useState(10);

  const fetchRoles = async () => {
    setLoading(true);
    try {
      const response = await request.post('/roles/page', {
        page: currentPage,
        pageSize: pageSize
      });
      setRoles(response.data);
      setTotal(response.total);
    } catch (error) {
      message.error('获取角色列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, [currentPage, pageSize]);

  const handleAdd = () => {
    form.resetFields();
    setEditingId(null);
    setModalVisible(true);
  };

  const handleEdit = async (record: Role) => {
    try {
      const response = await request.get(`/roles/detail/${record.id}`);
      form.setFieldsValue({
        name: response.name,
        code: response.code,
        description: response.description
      });
      setEditingId(record.id);
      setModalVisible(true);
    } catch (error) {
      message.error('获取角色详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个角色吗？',
      onOk: async () => {
        try {
          await request.delete(`/roles/${id}`);
          message.success('删除成功');
          fetchRoles();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const url = editingId ? `/roles/${editingId}` : '/roles';
      const method = editingId ? 'PUT' : 'POST';

      await request({
        url,
        method,
        data: values,
      });
      message.success(`${editingId ? '更新' : '创建'}成功`);
      setModalVisible(false);
      fetchRoles();
    } catch (error) {
      message.error('表单验证失败');
    }
  };

  const fetchPermissionList = async (page = 1) => {
    try {
      const response = await request.post('/permissions/page', {
        page,
        pageSize: permissionPageSize
      });
      setPermissionList(response.data);
      setPermissionTotal(response.total);
    } catch (error) {
      message.error('获取权限列表失败');
    }
  };

  const handlePermissionAssign = async (roleId: number) => {
    setCurrentRoleId(roleId);
    try {
      const response = await request.get(`/roles/${roleId}/permissions`);
      setPermissionList(response);
      setPermissionModalVisible(true);
    } catch (error) {
      message.error('获取角色权限失败');
    }
  };

  const handlePermissionSubmit = async () => {
    if (!currentRoleId) return;

    try {
      const selectedIds = getSelectedPermissionIds(permissionList);
      await request.post(`/roles/${currentRoleId}/permissions`, {
        permissions: selectedIds
      });
      message.success('权限分配成功');
      setPermissionModalVisible(false);
    } catch (error) {
      message.error('权限分配失败');
    }
  };

  const getSelectedPermissionIds = (permissions: Permission[]): string[] => {
    const selectedIds: string[] = [];
    const traverse = (items: Permission[]) => {
      items.forEach(item => {
        if (item.enable) {
          selectedIds.push(item.id.toString());
        }
        if (item.children && item.children.length > 0) {
          traverse(item.children);
        }
      });
    };
    traverse(permissions);
    return selectedIds;
  };

  const columns: ColumnsType<Role> = [
    {
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '角色编码',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
      render: (text: string) => new Date(text).toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      })
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="text"
            icon={<KeyOutlined />}
            onClick={() => handlePermissionAssign(record.id)}
          >
            分配权限
          </Button>
          <Button
            type="text"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增角色
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={roles}
        rowKey="id"
        loading={loading}
        pagination={{
          current: currentPage,
          pageSize: pageSize,
          total: total,
          onChange: (page, size) => {
            setCurrentPage(page);
            setPageSize(size);
          }
        }}
      />
      <Modal
        title={editingId ? '编辑角色' : '新增角色'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="code"
            label="角色编码"
            rules={[{ required: true, message: '请输入角色编码' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="分配权限"
        open={permissionModalVisible}
        onOk={handlePermissionSubmit}
        onCancel={() => setPermissionModalVisible(false)}
        width={800}
      >
        <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
          <Collapse>
            {permissionList && permissionList.length > 0 ? (
              permissionList.map(group => (
                <Collapse.Panel key={group.id} header={group.name}>
                  <Row gutter={[16, 16]}>
                    {group.children.map(item => (
                      <Col span={8} key={item.id}>
                        <Checkbox
                          checked={item.enable}
                          onChange={e => {
                            const { checked } = e.target;
                            const updatePermissionStatus = (items: Permission[], id: number): Permission[] => {
                              return items.map(item => {
                                if (item.id === id) {
                                  return { ...item, enable: checked };
                                }
                                if (item.children) {
                                  return {
                                    ...item,
                                    children: updatePermissionStatus(item.children, id)
                                  };
                                }
                                return item;
                              });
                            };
                            setPermissionList(prev => updatePermissionStatus(prev, item.id));
                          }}
                        >
                          {item.name}
                        </Checkbox>
                      </Col>
                    ))}
                  </Row>
                </Collapse.Panel>
              ))
            ) : (
              <p>没有可分配的权限</p>
            )}
          </Collapse>
        </div>
      </Modal>
    </div>
  );
};

export default RoleList;