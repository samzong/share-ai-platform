import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Input, 
  Row, 
  Col, 
  Tag, 
  Space, 
  Button,
  Typography,
  Select,
  message,
  Spin
} from 'antd';
import { SearchOutlined, StarOutlined, StarFilled } from '@ant-design/icons';
import { imageApi } from '../services/api';
import { ContainerImage, Label } from '../types/image';

const { Title, Paragraph } = Typography;
const { Search } = Input;

const Images: React.FC = () => {
  const [searchText, setSearchText] = useState('');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [images, setImages] = useState<ContainerImage[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);

  const fetchImages = async () => {
    try {
      setLoading(true);
      const response = await imageApi.getImages({
        page,
        page_size: pageSize,
        search: searchText,
      });
      console.log('API Response:', response);
      if (Array.isArray(response)) {
        setImages(response);
        setTotal(response.length);
      } else if (response && Array.isArray(response.data)) {
        setImages(response.data);
        setTotal(response.total || response.data.length);
      } else {
        setImages([]);
        setTotal(0);
      }
    } catch (error) {
      console.error('Error fetching images:', error);
      message.error('获取容器镜像列表失败');
      setImages([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchImages();
  }, [page, searchText]);

  const handleSearch = (value: string) => {
    setSearchText(value);
    setPage(1); // 重置页码
  };

  const handleTagChange = (value: string[]) => {
    setSelectedTags(value);
  };

  const handleCollect = async (imageId: string, isCollected: boolean) => {
    try {
      if (isCollected) {
        await imageApi.uncollectImage(imageId);
      } else {
        await imageApi.collectImage(imageId);
      }
      fetchImages(); // 刷新列表
      message.success(isCollected ? '取消收藏成功' : '收藏成功');
    } catch (error) {
      message.error('操作失败');
    }
  };

  // 获取所有标签
  const allTags = Array.from(
    new Set(images.flatMap(image => image.labels.map(label => label.name)))
  );

  const filteredImages = images.filter(image => {
    const matchesTags = selectedTags.length === 0 || 
                       selectedTags.every(tag => image.labels.some(label => label.name === tag));
    return matchesTags;
  });

  return (
    <div>
      <Title level={2}>容器镜像库</Title>
      <Space direction="vertical" size="middle" style={{ width: '100%', marginBottom: 24 }}>
        <Row gutter={16}>
          <Col span={12}>
            <Search
              placeholder="搜索容器镜像..."
              allowClear
              enterButton={<SearchOutlined />}
              size="large"
              onSearch={handleSearch}
            />
          </Col>
          <Col span={12}>
            <Select
              mode="multiple"
              style={{ width: '100%' }}
              placeholder="选择标签筛选"
              onChange={handleTagChange}
              options={allTags.map(tag => ({ label: tag, value: tag }))}
              size="large"
            />
          </Col>
        </Row>
      </Space>

      <Spin spinning={loading}>
        <Row gutter={[16, 16]}>
          {filteredImages.map(image => (
            <Col xs={24} sm={12} md={8} key={image.id}>
              <Card
                hoverable
                actions={[
                  <Button
                    type="text"
                    icon={image.stars > 0 ? <StarFilled style={{ color: '#faad14' }} /> : <StarOutlined />}
                    onClick={() => handleCollect(image.id, image.stars > 0)}
                  >
                    {image.stars}
                  </Button>,
                  <Button type="link">部署</Button>
                ]}
              >
                <Card.Meta
                  title={image.name}
                  description={
                    <>
                      <Paragraph>{image.description}</Paragraph>
                      <Space size={[0, 8]} wrap>
                        {image.labels.map(label => (
                          <Tag key={label.id} color="blue">{label.name}</Tag>
                        ))}
                      </Space>
                      <div style={{ marginTop: 8 }}>
                        <Space>
                          <span>{image.registry}/{image.namespace}/{image.repository}:{image.tag}</span>
                        </Space>
                      </div>
                    </>
                  }
                />
              </Card>
            </Col>
          ))}
        </Row>
      </Spin>
    </div>
  );
};

export default Images; 