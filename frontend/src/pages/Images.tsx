import React, { useState } from 'react';
import { 
  Card, 
  Input, 
  Row, 
  Col, 
  Tag, 
  Space, 
  Button,
  Typography,
  Select
} from 'antd';
import { SearchOutlined, StarOutlined, StarFilled } from '@ant-design/icons';

const { Title, Paragraph } = Typography;
const { Search } = Input;

interface ImageItem {
  id: string;
  name: string;
  description: string;
  tags: string[];
  stars: number;
  isStarred: boolean;
}

const Images: React.FC = () => {
  const [searchText, setSearchText] = useState('');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);

  // 模拟数据
  const mockImages: ImageItem[] = [
    {
      id: '1',
      name: 'Stable Diffusion',
      description: '开源的文本生成图像模型，支持多种风格和场景。',
      tags: ['AI', '图像生成', 'Stable Diffusion'],
      stars: 1200,
      isStarred: false,
    },
    {
      id: '2',
      name: 'ChatGLM',
      description: '开源的中文对话语言模型，支持多轮对话和知识问答。',
      tags: ['AI', '语言模型', 'ChatGLM'],
      stars: 980,
      isStarred: true,
    },
    {
      id: '3',
      name: 'YOLOv8',
      description: '实时目标检测模型，支持多种目标检测场景。',
      tags: ['AI', '目标检测', 'YOLO'],
      stars: 850,
      isStarred: false,
    },
  ];

  const allTags = Array.from(
    new Set(mockImages.flatMap(image => image.tags))
  );

  const handleSearch = (value: string) => {
    setSearchText(value);
  };

  const handleTagChange = (value: string[]) => {
    setSelectedTags(value);
  };

  const filteredImages = mockImages.filter(image => {
    const matchesSearch = image.name.toLowerCase().includes(searchText.toLowerCase()) ||
                         image.description.toLowerCase().includes(searchText.toLowerCase());
    const matchesTags = selectedTags.length === 0 || 
                       selectedTags.every(tag => image.tags.includes(tag));
    return matchesSearch && matchesTags;
  });

  return (
    <div>
      <Title level={2}>镜像库</Title>
      <Space direction="vertical" size="middle" style={{ width: '100%', marginBottom: 24 }}>
        <Row gutter={16}>
          <Col span={12}>
            <Search
              placeholder="搜索镜像..."
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

      <Row gutter={[16, 16]}>
        {filteredImages.map(image => (
          <Col xs={24} sm={12} md={8} key={image.id}>
            <Card
              hoverable
              actions={[
                image.isStarred ? 
                  <StarFilled style={{ color: '#faad14' }} /> : 
                  <StarOutlined />,
                <Button type="link">部署</Button>
              ]}
            >
              <Card.Meta
                title={image.name}
                description={
                  <>
                    <Paragraph>{image.description}</Paragraph>
                    <Space size={[0, 8]} wrap>
                      {image.tags.map(tag => (
                        <Tag key={tag} color="blue">{tag}</Tag>
                      ))}
                    </Space>
                    <div style={{ marginTop: 8 }}>
                      <StarOutlined /> {image.stars}
                    </div>
                  </>
                }
              />
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default Images; 