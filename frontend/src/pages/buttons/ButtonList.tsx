import { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, Form, Input, Select, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import request from '@/utils/request';
import type { ColumnsType } from 'antd/es/table';

interface SelectOption {
  label: string;
  value: number;
}

interface Button {
  id: number;
  name: string;
  code: string;
  menuId: number;
  createdAt: string;
}

const ButtonList = () => {
  const [buttons, setButtons] = useState<Button[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingId, setEditingId] = useState<number | null>(null);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [menuOptions, setMenuOptions] = useState<SelectOption[]>([]);
  const [menuLoading, setMenuLoading] = useState(false);

  const fetchButtons = async () => {
    setLoading(true);
    try {
      const response = await request.post('/buttons/page', {
        page: currentPage,
        pageSize: pageSize
      });
      setButtons(response.data);
      setTotal(response.total);
    } catch (error) {
      message.error('获取按钮列表失败');
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

  useEffect(() => {
    fetchButtons();
  }, [currentPage, pageSize]);

  useEffect(() => {
    if (modalVisible) {
      fetchMenuOptions();
    }
  }, [modalVisible]);

  const handleAdd = () => {
    form.resetFields();
    setEditingId(null);
    setModalVisible(true);
  };

  const handleEdit = async (record: Button) => {
    try {
      const response = await request.get(`/buttons/detail/${record.id}`);
      form.setFieldsValue({
        name: response.name,
        code: response.code,
        menuId: response.menu_id
      });
      setEditingId(record.id);
      setModalVisible(true);
    } catch (error) {
      message.error('获取按钮详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个按钮吗？',
      onOk: async () => {
        try {
          await request.delete(`/buttons/${id}`);
          message.success('删除成功');
          fetchButtons();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const url = editingId ? `/buttons/${editingId}` : '/buttons';
      const method = editingId ? 'PUT' : 'POST';

      await request({
        url,
        method,
        data: values,
      });
      message.success(`${editingId ? '更新' : '创建'}成功`);
      setModalVisible(false);
      fetchButtons();
    } catch (error) {
      message.error('表单验证失败');
    }
  };

  const columns: ColumnsType<Button> = [
    {
      title: '按钮名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '按钮编码',
      dataIndex: 'permission_code',
      key: 'permission_code',
    },
    {
      title: '所属菜单ID',
      dataIndex: 'menu_id',
      key: 'menu_id',
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
          新增按钮
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={buttons}
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
        title={editingId ? '编辑按钮' : '新增按钮'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="按钮名称"
            rules={[{ required: true, message: '请输入按钮名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="code"
            label="按钮编码"
            rules={[{ required: true, message: '请输入按钮编码' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="menuId"
            label="所属菜单"
            rules={[{ required: true, message: '请选择所属菜单' }]}
          >
            <Select
              showSearch
              loading={menuLoading}
              options={menuOptions}
              onSearch={(value) => fetchMenuOptions(value)}
              filterOption={false}
              placeholder="请选择所属菜单"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ButtonList;