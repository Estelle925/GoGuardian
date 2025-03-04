import { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, Form, Input, Select, message, InputNumber } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import request from '@/utils/request';
import type { ColumnsType } from 'antd/es/table';

interface SelectOption {
  label: string;
  value: number;
}

interface Menu {
  id: number;
  name: string;
  path: string;
  icon?: string;
  parentId?: number;
  order: number;
  createdAt: string;
}

const MenuList = () => {
  const [menus, setMenus] = useState<Menu[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [parentOptions, setParentOptions] = useState<SelectOption[]>([]);
  const [parentLoading, setParentLoading] = useState(false);

  const fetchMenus = async () => {
    setLoading(true);
    try {
      const response = await request.post('/menus/page', {
        page: currentPage,
        pageSize: pageSize
      });
      setMenus(response.data);
      setTotal(response.total);
    } catch (error) {
      message.error('获取菜单列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchParentOptions = async (search = '') => {
    setParentLoading(true);
    try {
      const response = await request.post('/menus/page', {
        page: 1,
        pageSize: 50,
        search
      });
      const options = response.data.map((menu: Menu) => ({
        label: menu.name,
        value: menu.id
      }));
      setParentOptions(options);
    } catch (error) {
      message.error('获取父级菜单列表失败');
    } finally {
      setParentLoading(false);
    }
  };

  useEffect(() => {
    fetchMenus();
  }, [currentPage, pageSize]);

  useEffect(() => {
    if (modalVisible) {
      fetchParentOptions();
    }
  }, [modalVisible]);

  const handleAdd = () => {
    form.resetFields();
    setEditingId(null);
    setModalVisible(true);
  };

  const handleEdit = async (record: Menu) => {
    try {
      const response = await request.get(`/menus/detail/${record.id}`);
      form.setFieldsValue({
        name: response.name,
        path: response.path,
        icon: response.icon,
        parentId: response.parent_id,
        order: response.order
      });
      setEditingId(record.id);
      setModalVisible(true);
    } catch (error) {
      message.error('获取菜单详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个菜单吗？',
      onOk: async () => {
        try {
          await request.delete(`/menus/${id}`);
          message.success('删除成功');
          fetchMenus();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const url = editingId ? `/menus/${editingId}` : '/menus';
      const method = editingId ? 'PUT' : 'POST';

      await request({
        url,
        method,
        data: values,
      });
      message.success(`${editingId ? '更新' : '创建'}成功`);
      setModalVisible(false);
      fetchMenus();
    } catch (error) {
      message.error('表单验证失败');
    }
  };

  const columns: ColumnsType<Menu> = [
    {
      title: '菜单名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '路径',
      dataIndex: 'path',
      key: 'path',
    },
    {
      title: '图标',
      dataIndex: 'icon',
      key: 'icon',
    },
    {
      title: '排序',
      dataIndex: 'order',
      key: 'order',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
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
          新增菜单
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={menus}
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
        title={editingId ? '编辑菜单' : '新增菜单'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="菜单名称"
            rules={[{ required: true, message: '请输入菜单名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="path"
            label="路径"
            rules={[{ required: true, message: '请输入路径' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="icon" label="图标">
            <Input />
          </Form.Item>
          <Form.Item name="parentId" label="父级菜单">
            <Select
              allowClear
              showSearch
              loading={parentLoading}
              options={parentOptions}
              onSearch={(value) => fetchParentOptions(value)}
              filterOption={false}
              placeholder="请选择父级菜单"
            />
          </Form.Item>
          <Form.Item
            name="order"
            label="排序"
            rules={[{ required: true, message: '请输入排序序号' }]}
          >
            <InputNumber min={0} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default MenuList;