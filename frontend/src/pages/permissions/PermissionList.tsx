import { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, Form, Input, Select, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import request from '@/utils/request';
import type { ColumnsType } from 'antd/es/table';

interface Permission {
  id: number;
  name: string;
  permissionCode: string;
  type: string;
  menuId?: number;
  buttonId?: number;
  parentId?: number;
  createdAt: string;
}

interface SelectOption {
  label: string;
  value: number;
}

const PermissionList = () => {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [menuOptions, setMenuOptions] = useState<SelectOption[]>([]);
  const [buttonOptions, setButtonOptions] = useState<SelectOption[]>([]);
  const [parentOptions, setParentOptions] = useState<SelectOption[]>([]);
  const [menuLoading, setMenuLoading] = useState(false);
  const [buttonLoading, setButtonLoading] = useState(false);
  const [parentLoading, setParentLoading] = useState(false);
  const [permissionType, setPermissionType] = useState<string>('menu');

  // 监听权限类型变化
  const handleTypeChange = (value: string) => {
    setPermissionType(value);
    form.setFieldsValue({
      menuId: undefined,
      buttonId: undefined
    });
  };

  const fetchPermissions = async () => {
    setLoading(true);
    try {
      const response = await request.post('/permissions/page', {
        page: currentPage,
        pageSize: pageSize
      });
      setPermissions(response.data);
      setTotal(response.total);
    } catch (error) {
      message.error('获取权限列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchMenuOptions = async (search = '') => {
    setMenuLoading(true);
    try {
      const response = await request.post('/menus/page', {
        page: 1,
        pageSize: 50,
        search
      });
      const options = response.data.map((menu: any) => ({
        label: menu.name,
        value: menu.id
      }));
      setMenuOptions(options);
    } catch (error) {
      message.error('获取菜单列表失败');
    } finally {
      setMenuLoading(false);
    }
  };

  const fetchButtonOptions = async (search = '') => {
    setButtonLoading(true);
    try {
      const response = await request.post('/buttons/page', {
        page: 1,
        pageSize: 50,
        search
      });
      const options = response.data.map((button: any) => ({
        label: button.name,
        value: button.id
      }));
      setButtonOptions(options);
    } catch (error) {
      message.error('获取按钮列表失败');
    } finally {
      setButtonLoading(false);
    }
  };

  const fetchParentOptions = async (search = '') => {
    setParentLoading(true);
    try {
      const response = await request.post('/permissions/page', {
        page: 1,
        pageSize: 50,
        search
      });
      const options = response.data.map((permission: Permission) => ({
        label: permission.name,
        value: permission.id
      }));
      setParentOptions(options);
    } catch (error) {
      message.error('获取父权限列表失败');
    } finally {
      setParentLoading(false);
    }
  };

  useEffect(() => {
    fetchPermissions();
  }, [currentPage, pageSize]);

  useEffect(() => {
    if (modalVisible) {
      fetchMenuOptions();
      fetchButtonOptions();
      fetchParentOptions();
    }
  }, [modalVisible]);

  const handleAdd = () => {
    form.resetFields();
    setEditingId(null);
    setModalVisible(true);
  };

  const handleEdit = async (record: Permission) => {
    try {
      const response = await request.get(`/permissions/detail/${record.id}`);
      // 将后端返回的数据字段映射到表单字段
      form.setFieldsValue({
        name: response.name,
        permissionCode: response.code,
        type: response.type,
        parentId: response.parent_id,
        menuId: response.menu_id,
        buttonId: response.button_id
      });
      setEditingId(record.id);
      setModalVisible(true);
    } catch (error) {
      message.error('获取权限详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个权限吗？',
      onOk: async () => {
        try {
          await request.delete(`/permissions/${id}`);
          message.success('删除成功');
          fetchPermissions();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const url = editingId ? `/permissions/${editingId}` : '/permissions';
      const method = editingId ? 'PUT' : 'POST';

      await request({
        url,
        method,
        data: values,
      });
      message.success(`${editingId ? '更新' : '创建'}成功`);
      setModalVisible(false);
      fetchPermissions();
    } catch (error) {
      message.error('表单验证失败');
    }
  };

  const columns: ColumnsType<Permission> = [
    {
      title: '权限名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '权限编码',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '权限类型',
      dataIndex: 'type',
      key: 'type',
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
          新增权限
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={permissions}
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
        title={editingId ? '编辑权限' : '新增权限'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="权限名称"
            rules={[{ required: true, message: '请输入权限名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="permissionCode"
            label="权限编码"
            rules={[{ required: true, message: '请输入权限编码' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="type"
            label="权限类型"
            rules={[{ required: true, message: '请选择权限类型' }]}
          >
            <Select onChange={handleTypeChange}>
              <Select.Option value="menu">菜单权限</Select.Option>
              <Select.Option value="button">按钮权限</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="parentId"
            label="上级权限"
            rules={[{ required: false }]}
          >
            <Select
              allowClear
              showSearch
              loading={parentLoading}
              options={parentOptions}
              onSearch={(value) => fetchParentOptions(value)}
              filterOption={false}
              placeholder="请选择上级权限"
            />
          </Form.Item>
          <Form.Item
            name="menuId"
            label="关联菜单"
            rules={[{ required: permissionType === 'menu', message: '请选择关联菜单' }]}
            style={{ display: permissionType === 'menu' ? 'block' : 'none' }}
          >
            <Select
              allowClear
              showSearch
              loading={menuLoading}
              options={menuOptions}
              onSearch={(value) => fetchMenuOptions(value)}
              filterOption={false}
              placeholder="请选择关联菜单"
            />
          </Form.Item>
          <Form.Item
            name="buttonId"
            label="关联按钮"
            rules={[{ required: permissionType === 'button', message: '请选择关联按钮' }]}
            style={{ display: permissionType === 'button' ? 'block' : 'none' }}
          >
            <Select
              allowClear
              showSearch
              loading={buttonLoading}
              options={buttonOptions}
              onSearch={(value) => fetchButtonOptions(value)}
              filterOption={false}
              placeholder="请选择关联按钮"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default PermissionList;