import React from 'react';
import { Typography, Card, Row, Col, Button } from 'antd';
import { useNavigate } from 'react-router-dom';

const { Title, Paragraph } = Typography;

const Home: React.FC = () => {
  const navigate = useNavigate();

  const features = [
    {
      title: '镜像管理',
      description: '浏览、搜索和收藏 AI 相关镜像，快速找到你需要的应用。',
      action: () => navigate('/images'),
    },
    {
      title: '一键部署',
      description: '选择合适的云平台，一键部署你的 AI 应用。',
      action: () => navigate('/images'),
    },
    {
      title: '开源共享',
      description: '加入我们的开源社区，分享你的 AI 应用镜像。',
      action: () => window.open('https://github.com/samzong/share-ai-platform', '_blank'),
    },
  ];

  return (
    <div style={{ padding: '24px 0' }}>
      <Typography style={{ textAlign: 'center', marginBottom: 48 }}>
        <Title>Share AI Platform</Title>
        <Paragraph>
          开源的 AI 镜像分享平台，致力于简化 AI 应用的部署和分享过程
        </Paragraph>
      </Typography>

      <Row gutter={[24, 24]} justify="center">
        {features.map((feature, index) => (
          <Col xs={24} sm={12} md={8} key={index}>
            <Card
              hoverable
              title={feature.title}
              style={{ height: '100%' }}
              actions={[
                <Button type="primary" onClick={feature.action}>
                  了解更多
                </Button>,
              ]}
            >
              <Paragraph>{feature.description}</Paragraph>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default Home; 