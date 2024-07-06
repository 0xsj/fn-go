import React, { PureComponent } from "react";
import {
  Radar,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  ResponsiveContainer,
} from "recharts";

const data = [
  {
    department: "Sales",
    Efficiency: 85,
    CustomerSatisfaction: 90,
    Revenue: 80,
    Innovation: 75,
    Compliance: 95,
  },
  {
    department: "Marketing",
    Efficiency: 80,
    CustomerSatisfaction: 85,
    Revenue: 70,
    Innovation: 85,
    Compliance: 90,
  },
  {
    department: "Operations",
    Efficiency: 90,
    CustomerSatisfaction: 85,
    Revenue: 75,
    Innovation: 80,
    Compliance: 85,
  },
  {
    department: "Finance",
    Efficiency: 95,
    CustomerSatisfaction: 80,
    Revenue: 85,
    Innovation: 70,
    Compliance: 90,
  },
  {
    department: "Legal",
    Efficiency: 85,
    CustomerSatisfaction: 90,
    Revenue: 80,
    Innovation: 75,
    Compliance: 95,
  },
];

export const EmployeePerformanceRadarChart = () => (
  <ResponsiveContainer width='100%' height={400}>
    <RadarChart cx='50%' cy='50%' outerRadius='80%' data={data}>
      <PolarGrid />
      <PolarAngleAxis dataKey='department' />
      <PolarRadiusAxis />
      <Radar
        name='Performance Metrics'
        dataKey='Efficiency'
        stroke='#8884d8'
        fill='#8884d8'
        fillOpacity={0.6}
      />
      <Radar
        name='Performance Metrics'
        dataKey='CustomerSatisfaction'
        stroke='#82ca9d'
        fill='#82ca9d'
        fillOpacity={0.6}
      />
      <Radar
        name='Performance Metrics'
        dataKey='Revenue'
        stroke='#ffc658'
        fill='#ffc658'
        fillOpacity={0.6}
      />
      <Radar
        name='Performance Metrics'
        dataKey='Innovation'
        stroke='#FF5733'
        fill='#FF5733'
        fillOpacity={0.6}
      />
      <Radar
        name='Performance Metrics'
        dataKey='Compliance'
        stroke='#6A5ACD'
        fill='#6A5ACD'
        fillOpacity={0.6}
      />
    </RadarChart>
  </ResponsiveContainer>
);
