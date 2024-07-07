import { AreaGradientChart } from "@/components/charts/area/area-gradient";
import { StackedAreaChart } from "@/components/charts/area/stacked";
import { StackedExpandedAreaChart } from "@/components/charts/area/stacked-expanded";
import { InteractiveBarChart } from "@/components/charts/bar/interactive";
import { InteractiveMixedChart } from "@/components/charts/bar/mixed";
import { InteractiveMultipleBarChart } from "@/components/charts/bar/multple";
import { NegativeBarChart } from "@/components/charts/bar/negative";
import { StackedBarChart } from "@/components/charts/bar/stacked";
import { MultipleLineChart } from "@/components/charts/line/multiple";
import { SingleLineChart } from "@/components/charts/line/single";
import { StepLineChart } from "@/components/charts/line/step";
import { PieDonutChart } from "@/components/charts/pie/donut";
import { DonutPieTextChart } from "@/components/charts/pie/donut-with-text";
import { PieWithLabelChart } from "@/components/charts/pie/label";
import { PieLegendChart } from "@/components/charts/pie/legend";
import { RadarLineOnlyChart } from "@/components/charts/radar/lines-only";
import { RadialGridChart } from "@/components/charts/radial/grid";
import ComingSoon from "@/components/coming-soon";
import { TabsContent } from "@/components/ui/tabs";

interface Props {}

export const ReportsTab: React.FC<Props> = () => {
  return (
    <TabsContent value='reports'>
      <div className='gap-6 md:flex md:flex-row-reverse md:items-start'>
        <StackedAreaChart />
        <StackedExpandedAreaChart />
        <AreaGradientChart />
        <InteractiveMixedChart />
        <InteractiveMultipleBarChart />
        <NegativeBarChart />
        <StackedBarChart />
      </div>
      <div className='gap-6 md:flex md:flex-row-reverse md:items-start'>
        <MultipleLineChart />
        <SingleLineChart />
        <StepLineChart />
      </div>
      <div className='gap-6 md:flex md:flex-row-reverse md:items-start'>
        <DonutPieTextChart />
        <PieDonutChart />
        <PieWithLabelChart />
        <PieLegendChart />
      </div>
      <div className='gap-6 md:flex md:flex-row-reverse md:items-start'>
        <RadarLineOnlyChart />
        <RadialGridChart />
      </div>
    </TabsContent>
  );
};
